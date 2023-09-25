package service

type User struct {
	ID   int64  `json:"ID"`
	Name string `json:"name"`
}

var NilUser User
