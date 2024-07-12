package handler

import (
	"broker-td/config"
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"testing"
)

const _gdfuPath = "/gdrive-family-uploader"

type RoundTripFunc func(hc *fasthttp.HostClient, req *fasthttp.Request, resp *fasthttp.Response) bool

func (f RoundTripFunc) RoundTrip(
	hc *fasthttp.HostClient,
	req *fasthttp.Request,
	resp *fasthttp.Response) (retry bool, err error) {
	return f(hc, req, resp), nil
}

func NewTestUserClient(fn RoundTripFunc) *fiber.Agent {
	agent := fiber.AcquireAgent()
	if err := agent.Parse(); err != nil {
		return nil
	}

	agent.Transport = fn

	return agent
}

// Server is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type Server struct {
}

func (s *Server) UploadDriveFile(payload FamilyPayload, resp *string) error {
	*resp = "success"
	return nil
}

func rpcListen(listen net.Listener) error {
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(listen)

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func startRPCServer() (string, error) {
	if err := rpc.Register(new(Server)); err != nil {
		return "", err
	}

	listen, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	address := listen.Addr().String()
	log.Println("Starting RPC Server on:", address)

	go func() {
		err := rpcListen(listen)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}()

	return address, nil
}

func createGDriveFamilyServer(userClient *fiber.Agent) *fiber.App {
	app := fiber.New()

	rpcAddress, err := startRPCServer()

	configs := &config.EnvVars{
		UserServiceBaseURL:  "localhost:3001",
		DriveServiceBaseRPC: rpcAddress,
	}

	gdfuHandler, err := NewGDriveFamilyHandler(configs)
	if err != nil {
		return nil
	}

	gdfuHandler.Client = userClient

	app.Route("/gdrive-family-uploader", func(api fiber.Router) {
		api.Post("/", gdfuHandler.Post).Name("post")
	}, "gdrive-family-uploader.")

	return app
}

func createGDriveFamilyRequest(
	method string, url string, buffer bytes.Buffer, writer *multipart.Writer) (*http.Request, error) {
	req := httptest.NewRequest(method, url, bytes.NewReader(buffer.Bytes()))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

func TestGDriveFamilyHandlerPost_Successful(t *testing.T) {
	// Given
	userClientResponseBody := `
		{
			"folder_id": "test"
		}
	`

	userClient := NewTestUserClient(func(hc *fasthttp.HostClient, req *fasthttp.Request, resp *fasthttp.Response) bool {
		resp.SetStatusCode(fiber.StatusOK)
		resp.SetBody([]byte(userClientResponseBody))

		return false
	})

	server := createGDriveFamilyServer(userClient)

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	_, err := writer.CreateFormFile("tg-bot-file", "test.jpg")
	if err != nil {
		require.Error(t, err)
	}

	err = writer.WriteField("username", "test")
	if err != nil {
		require.Error(t, err)
	}

	err = writer.WriteField("bot_name", "test")
	if err != nil {
		require.Error(t, err)
	}

	err = writer.WriteField("date", "test")
	if err != nil {
		require.Error(t, err)
	}

	// Send date in order to identify the folder.
	err = writer.WriteField("filename", "test.jpg")
	if err != nil {
		require.Error(t, err)
	}

	if err = writer.Close(); err != nil {
		require.Error(t, err)
	}

	req, err := createGDriveFamilyRequest(fiber.MethodPost, _gdfuPath, buffer, writer)
	require.NoError(t, err)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var data any
	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	require.NotEmpty(t, data)
}
