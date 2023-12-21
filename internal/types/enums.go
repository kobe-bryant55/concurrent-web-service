package types

type Status string

const (
	Active  Status = "active"
	Passive Status = "passive"
)

type Role string

var (
	Admin      Role = "admin"
	Registered Role = "registered"
)
