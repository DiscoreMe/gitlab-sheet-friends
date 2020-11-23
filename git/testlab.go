package git

import (
	"time"
)

type TestLab struct {
	issues []Issue
}

func (t *TestLab) CreateClient() error {
	t.issues = []Issue{
		NewFakeIssue(1, 1, "fake issue #1", "???", "opened", NewDate(6, 04, 2020)),
		NewFakeIssue(2, 2, "fake issue #2", "???", "closed", NewDate(7, 04, 2020)),
		NewFakeIssue(3, 3, "fake issue #3", "???", "closed", NewDate(7, 04, 2020)),
		NewFakeIssue(4, 4, "fake issue #4", "???", "opened", NewDate(8, 04, 2020)),
		NewFakeIssue(5, 5, "fake issue #5", "???", "opened", NewDate(8, 04, 2020)),
	}
	return nil
}
func (t *TestLab) Issues(IssuesOptions) ([]Issue, error) {
	return t.issues, nil
}

func (t *TestLab) IsClosed(state string) bool {
	if state == "closed" {
		return true
	}
	return false
}

func (t *TestLab) IssuesByIDs(IssuesOptions) ([]Issue, error) {
	return t.issues, nil
}

func (t *TestLab) Ping() error {
	return nil
}

func (t *TestLab) Availability() string {
	return AvailableExternal
}

func (t *TestLab) AddMember(string) {}

func (t *TestLab) SetConf(HostConf) {}

func NewFakeIssue(issueID, projectID int, title, link, state string, createdAt time.Time) Issue {
	return Issue{
		ID:        issueID,
		ProjectID: projectID,
		Title:     title,
		Link:      link,
		State:     state,
		CreatedAt: createdAt,
	}
}

func NewDate(d, m, y int) time.Time {
	return time.Date(y, time.Month(m), d, 22, 30, 15, 0, &time.Location{})
}
