package models

// Server is a database structure that keeps track of the individual server and if it is up
type Server struct {
	BaseModel
	DNS    string `json:"dns"`
	Region string `json:"region"`
	User   []User `json:"user"`
}
