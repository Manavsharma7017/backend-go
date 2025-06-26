package models

type Role string

const (
	RoleAdmin  Role = "ADMIN"
	RoleEditor Role = "EDITOR"
	RoleUSER   Role = "USER"
)
