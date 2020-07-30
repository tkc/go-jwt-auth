package models

type RoleType string

var RoleTypes = struct {
	Admin RoleType
	User  RoleType
}{
	Admin: "admin",
	User:  "user",
}

type User struct {
	ID   int
	Role RoleType
	Name string
	//Email string
}
