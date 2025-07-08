package models

import "github.com/google/uuid"

// User is a database structure that keeps track of each user connected to the VM currently.
type User struct {
	BaseModel
	Username string     `json:"username"`
	ServerID *uuid.UUID `json:"server_id"`
}
