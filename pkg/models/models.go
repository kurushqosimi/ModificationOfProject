package models

import "time"

type Config struct {
	ServerSetting ServerSettings `json:"server_setting"`
	DBSetting     DBSettings     `json:"db_setting"`
}

type DBSettings struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type ServerSettings struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Notes struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
	UserID    int       `json:"user_id"`
	Active    bool      `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
