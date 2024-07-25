package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"user-service/config"
	"user-service/internal/botuser"
	"user-service/internal/domain"
	"user-service/internal/platform/client"
)

const _defaultFolderCreatorCronSpec = "@yearly"

type FolderCronJob struct {
	ctx            context.Context
	configs        *config.EnvVars
	cron           *cron.Cron
	botUserService botuser.Service
	brokerClient   *client.BrokerClient
	logger         *zap.Logger
}

func NewFolderCronJob(
	ctx context.Context,
	configs *config.EnvVars,
	cron *cron.Cron,
	botUserService botuser.Service,
	brokerClient *client.BrokerClient,
	logger *zap.Logger,
) *FolderCronJob {
	return &FolderCronJob{
		ctx:            ctx,
		configs:        configs,
		cron:           cron,
		botUserService: botUserService,
		brokerClient:   brokerClient,
		logger:         logger,
	}
}

func (f *FolderCronJob) Run() error {
	cronSpec := _defaultFolderCreatorCronSpec
	if f.configs.IsDevelopment {
		cronSpec = "* * * * *"
	}

	entryID, err := f.cron.AddFunc(cronSpec, f.createYearlyFolders)
	if err != nil {
		return err
	}

	f.logger.Info(fmt.Sprintf("Folder cron job added, ID: %d", entryID))

	f.cron.Start()

	return nil
}

func (f *FolderCronJob) createYearlyFolders() {
	botUsers, err := f.botUserService.GetAllParents(f.ctx)
	if err != nil {
		f.logger.Error(err.Error())
		return
	}

	b, err := json.Marshal(botUsers)
	if err != nil {
		f.logger.Error(err.Error())
		return
	}

	response, err := f.brokerClient.PostFolderParentsCreations(b)
	if err != nil {
		return
	}

	data, ok := response["data"]
	if !ok {
		f.logger.Error("cannot process bot user empty")
		return
	}

	body, err := json.Marshal(data)
	if err != nil {
		f.logger.Error(err.Error())
		return
	}

	var newBotUserChildren []domain.BotUser
	if err := json.Unmarshal(body, &newBotUserChildren); err != nil {
		f.logger.Error(err.Error())
		return
	}

	if err := f.botUserService.SaveMany(f.ctx, newBotUserChildren); err != nil {
		f.logger.Error(err.Error())
		return
	}
}
