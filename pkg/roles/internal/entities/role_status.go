package entities

type RoleStatus uint8

const (
	Active  RoleStatus = 0
	Pending RoleStatus = 1
	Banned  RoleStatus = 2
	UnFound RoleStatus = 3
)
