package domain

type User struct {
	ID       int    `json:"id" db:"id" `
	Username string `json:"username" db:"username"`
}
