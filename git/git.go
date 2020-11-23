package git

import "time"

const (
	AvailableInternal = "internal"
	AvailableExternal = "external"
)

type Host interface {
	CreateClient() error
	Ping() error
	SetConf(HostConf)
	Issues(IssuesOptions) ([]Issue, error)
	IsClosed(string) bool
	IssuesByIDs(IssuesOptions) ([]Issue, error)
	Availability() string
	AddMember(string)
}

type HostConf struct {
	Type      string
	URL       string
	Available string
	Token     string
	Members   []string
}

type Issue struct {
	ID        int
	ProjectID int
	Title     string
	Link      string
	CreatedAt time.Time
	State     string
}

type IssuesOptions struct {
	State        string
	Scope        string
	CreatedAfter time.Time
	Sort         string
	IDs          []int
	Available    string
}
