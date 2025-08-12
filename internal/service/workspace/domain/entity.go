package domain

type Workspace struct {
	ID      uint
	Name    string
	Owner   Owner
	Members []Member
}

type Owner struct {
	ID       uint
	Username string
}

type Member struct {
	ID       uint
	Username string
	Access   Access
}

type Access string

const (
	AccessWrite Access = "workspace:write"
	AccessRead  Access = "workspace:read"
)
