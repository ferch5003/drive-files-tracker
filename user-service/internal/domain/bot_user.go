package domain

type BotUser struct {
	BotID              int    `json:"bot_id" db:"bot_id"`
	UserID             int    `json:"user_id" db:"user_id"`
	Date               string `json:"date" db:"date"`
	FolderID           string `json:"folder_id" db:"folder_id"`
	IsParent           bool   `json:"is_parent" db:"is_parent"`
	SpreadsheetID      string `json:"spreadsheet_id" db:"spreadsheet_id"`
	SpreadsheetGID     string `json:"spreadsheet_gid" db:"spreadsheet_gid"`
	SpreadsheetBaseGID string `json:"spreadsheet_base_gid" db:"spreadsheet_base_gid"`
	SpreadsheetColumn  string `json:"spreadsheet_column"  db:"spreadsheet_column"`
}
