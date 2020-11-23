package git

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/xanzy/go-gitlab"
	"sort"
)

type GitLab struct {
	client *gitlab.Client
	HostConf
}

func (g *GitLab) CreateClient() (err error) {
	client, err := gitlab.NewClient(g.Token, gitlab.WithBaseURL(g.URL))
	g.client = client
	return
}

func (g *GitLab) Ping() error {
	_, _, err := g.client.Issues.ListIssues(&gitlab.ListIssuesOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (g *GitLab) SetConf(conf HostConf) {
	g.Token = conf.Token
	g.Available = conf.Available
	g.URL = conf.URL
	g.Type = conf.Type
}

func (g *GitLab) Issues(options IssuesOptions) ([]Issue, error) {
	opts := &gitlab.ListIssuesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	}
	if options.Scope != "" {
		opts.Scope = &options.Scope
	}
	if options.State != "" {
		opts.State = &options.State
	}
	if options.Sort != "" {
		opts.Sort = &options.Sort
	}

	if !options.CreatedAfter.IsZero() {
		opts.CreatedAfter = &options.CreatedAfter
	}

	issues, err := func() ([]*gitlab.Issue, error) {
		if options.Available != "" {
			if options.Available == AvailableInternal {
				return g.internalIssues(opts)
			}
		}
		return g.externalIssues(opts)
	}()
	if err != nil {
		return nil, err
	}

	var iss []Issue
	for _, issue := range issues {
		iss = append(iss, Issue{
			ID:        issue.IID,
			ProjectID: issue.ProjectID,
			Title:     issue.Title,
			Link:      issue.WebURL,
			CreatedAt: *issue.CreatedAt,
			State:     issue.State,
		})
	}

	return iss, nil
}

func (g *GitLab) externalIssues(opts *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	issues, _, err := g.client.Issues.ListIssues(opts)
	return issues, err
}

func (g *GitLab) internalIssues(opts *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	if len(g.Members) == 0 {
		return []*gitlab.Issue{}, nil
	}
	var issues []*gitlab.Issue
	for _, member := range g.Members {
		iss, _, err := g.client.Issues.ListIssues(opts, func(req *retryablehttp.Request) error {
			q := req.URL.Query()
			q.Add("assignee_username", member)
			req.URL.RawQuery = q.Encode()
			return nil
		})
		if err != nil {
			return nil, err
		}
		issues = append(issues, iss...)
	}

	for _, member := range g.Members {
		iss, _, err := g.client.Issues.ListIssues(opts, func(req *retryablehttp.Request) error {
			q := req.URL.Query()
			q.Add("author_username", member)
			req.URL.RawQuery = q.Encode()
			return nil
		})
		if err != nil {
			return nil, err
		}
		issues = append(issues, iss...)
	}

	if len(issues) == 0 {
		return issues, nil
	}

	var processedIssues []*gitlab.Issue
	for i := 0; i < len(issues); i++ {
		var exist bool
		for j := 0; j < len(processedIssues); j++ {
			if issues[i].ID == processedIssues[j].ID {
				exist = true
				break
			}
		}
		if !exist {
			processedIssues = append(processedIssues, issues...)
		}
	}
	sort.Slice(processedIssues, func(i, j int) bool {
		return processedIssues[i].ID < processedIssues[j].ID
	})

	return processedIssues, nil
}

func (g *GitLab) IssuesByIDs(options IssuesOptions) ([]Issue, error) {
	opts := &gitlab.ListIssuesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
		IIDs:        options.IDs,
	}
	if options.Scope != "" {
		opts.Scope = &options.Scope
	}
	if options.State != "" {
		opts.State = &options.State
	}
	if options.Sort != "" {
		opts.Sort = &options.Sort
	}
	if !options.CreatedAfter.IsZero() {
		opts.CreatedAfter = &options.CreatedAfter
	}

	issues, _, err := g.client.Issues.ListIssues(opts)
	if err != nil {
		return nil, err
	}
	var iss []Issue
	for _, issue := range issues {
		iss = append(iss, Issue{
			ID:        issue.IID,
			ProjectID: issue.ProjectID,
			Title:     issue.Title,
			Link:      issue.WebURL,
			CreatedAt: *issue.CreatedAt,
			State:     issue.State,
		})
	}

	return iss, nil
}

func (g *GitLab) AddMember(username string) {
	g.Members = append(g.Members, username)
}

func (g *GitLab) IsClosed(state string) bool {
	if state == "closed" {
		return true
	}
	return false
}

func (g *GitLab) Availability() string {
	return g.Available
}
