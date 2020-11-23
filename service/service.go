package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/DiscoreMe/gitlab-sheets-friends/config"
	"github.com/DiscoreMe/gitlab-sheets-friends/git"
	"github.com/DiscoreMe/gitlab-sheets-friends/sheets"
	"github.com/DiscoreMe/gitlab-sheets-friends/storage"
)

const startRangeRow = 3

type Service struct {
	cfg  *config.Config
	stor *storage.Storage
	gits map[string]git.Host
	svr  *sheets.Sheets

	incompleteIssuesHandler func() error
}

func NewService(cfg *config.Config, stor *storage.Storage) (*Service, error) {
	var gits, err = initHosts(cfg.Gits)
	if err != nil {
		return nil, err
	}
	cfg.Gits = nil
	gits = initMembers(gits, cfg.Members)

	svr, err := sheets.NewSheets(cfg.SpreadSheetID)
	if err != nil {
		if _, ok := err.(sheets.ErrTokenNotFoundOrOutdated); !ok {
			return nil, err
		}
	}

	return &Service{
		cfg:  cfg,
		stor: stor,
		gits: gits,
		svr:  svr,
	}, err
}

func initHosts(cfgGits []config.Git) (map[string]git.Host, error) {
	var gits = make(map[string]git.Host)
	for _, g := range cfgGits {
		var host git.Host
		switch g.Type {
		case "gitlab":
			host = &git.GitLab{}
		case "testlab":
			host = &git.TestLab{}
		default:
			return nil, fmt.Errorf("error %s: type %s is unknown", g.Name, g.Type)
		}
		host.SetConf(git.HostConf{
			Type:      g.Type,
			URL:       g.URL,
			Available: g.Available,
			Token:     g.Token,
		})
		if err := host.CreateClient(); err != nil {
			return nil, err
		}
		gits[g.Name] = host
	}
	return gits, nil
}

func initMembers(gits map[string]git.Host, members []config.Member) map[string]git.Host {
	for _, member := range members {
		for key, value := range member.Services {
			if _, ok := gits[key]; !ok {
				fmt.Printf("warning: member %s: control version system %s is not found", member.ID, key)
				continue
			}
			gits[key].AddMember(value)
		}
	}
	return gits
}

func (s *Service) Run() error {
	var lists = make(map[string]*sheets.Sender)
	var counts = make(map[string]int)

	s.defaultIncompleteIssuesHandler()

	for key, cvs := range s.gits {
		createdAfter, err := s.createdAfterTimeIssue(key)
		if err != nil {
			return err
		}

		options := git.IssuesOptions{
			State:        "all",
			Scope:        "all",
			Sort:         "asc",
			CreatedAfter: createdAfter,
			Available:    cvs.Availability(),
		}

		issues, err := cvs.Issues(options)
		if err != nil {
			return err
		}
		for _, issue := range issues {
			listName := TimeListName(issue.CreatedAt)
			if _, ok := lists[listName]; !ok {
				lists[listName] = &sheets.Sender{}
				counts[listName] = 0
			}

			if err := s.sheetChecker(listName); err != nil {
				return err
			}

			listID, err := s.stor.ListIDByName(listName)
			if err != nil {
				return err
			}

			_, err = s.stor.IssueByGitProjectIDIssueIDListID(key, issue.ProjectID, issue.ID, listID)
			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}
				if err := s.stor.InsertIssue(&storage.Issue{
					Git:       key,
					ProjectID: issue.ProjectID,
					IssueID:   issue.ID,
					CreatedAt: issue.CreatedAt,
					ListID:    listID,
					IsClosed:  cvs.IsClosed(issue.State),
				}); err != nil {
					return err
				}

				lists[listName].AddRows(key, issue.Title, issue.Link, issue.CreatedAt)
				counts[listName]++
			}
		}
	}

	for name, sender := range lists {
		listID, err := s.stor.ListIDByName(name)
		if err != nil {
			return err
		}
		count, err := s.stor.CountIssuesByListID(listID)
		if err != nil {
			return err
		}

		startRow := count - counts[name]
		if startRow < 0 {
			log.Printf("warning: startRow %s is equal to %d\n", name, startRow)
			startRow = 2
		}

		sender.SetStartRange("A", startRangeRow+startRow)

		if err := s.svr.Update(sender, name); err != nil {
			return err
		}
	}

	if err := s.incompleteIssuesHandler(); err != nil {
		return err
	}

	return nil
}

func (s *Service) defaultIncompleteIssuesHandler() {
	listID := func(t time.Time) (int, error) {
		listName := TimeListName(t)
		listID, err := s.stor.ListIDByName(listName)
		if err != nil {
			return 0, err
		}
		return listID, nil
	}

	s.incompleteIssuesHandler = func() error {
		unclosedIssues, err := s.stor.IssuesByIsClosed(false)
		if err != nil {
			return err
		}

		var closedIssues []storage.Issue
		var openedIssuesGitIDs = make(map[string][]int)
		var overdueIssues []storage.Issue

		for _, issue := range unclosedIssues {
			openedIssuesGitIDs[issue.Git] = append(openedIssuesGitIDs[issue.Git], issue.ID)
		}
		for key, value := range openedIssuesGitIDs {
			createdAfter, err := s.createdAfterTimeIssue(key)
			if err != nil {
				return err
			}

			iss, err := s.gits[key].IssuesByIDs(git.IssuesOptions{
				State:        "all",
				Scope:        "all",
				CreatedAfter: createdAfter,
				Sort:         "asc",
				IDs:          value,
			})
			if err != nil {
				return err
			}
			for _, is := range iss {
				listID, err := listID(is.CreatedAt)
				if err != nil {
					return err
				}
				sis := storage.Issue{
					Git:       key,
					ProjectID: is.ProjectID,
					IssueID:   is.ID,
					CreatedAt: is.CreatedAt,
					ListID:    listID,
					IsClosed:  s.gits[key].IsClosed(is.State),
				}
				if s.gits[key].IsClosed(is.State) {
					closedIssues = append(closedIssues, sis)
				} else {
					overdueIssues = append(overdueIssues, sis)
				}
			}
		}

		nowLabel := TimeListName(time.Now())
		var processedOverdueIssues []storage.Issue
		for _, issue := range overdueIssues {
			issueLabel := TimeListName(issue.CreatedAt)
			if issueLabel != nowLabel {
				var is = issue
				is.CreatedAt = time.Now()
				issue.IsClosed = true
				if err := s.sheetChecker(issueLabel); err != nil {
					return err
				}
				if err := s.sheetChecker(nowLabel); err != nil {
					return err
				}
				processedOverdueIssues = append(processedOverdueIssues, is)
				closedIssues = append(closedIssues, issue)
			}
		}

		var sender = &sheets.Sender{}
		var count int
		var issuesIDs = make(map[string][]int)

		for _, issue := range processedOverdueIssues {
			issuesIDs[issue.Git] = append(issuesIDs[issue.Git], issue.IssueID)
		}

		for key, value := range issuesIDs {
			createdAfter, err := s.createdAfterTimeIssue(key)
			if err != nil {
				return err
			}

			iss, err := s.gits[key].IssuesByIDs(git.IssuesOptions{
				State:        "all",
				Scope:        "all",
				CreatedAfter: createdAfter,
				Sort:         "asc",
				IDs:          value,
				Available:    s.gits[key].Availability(),
			})
			if err != nil {
				return err
			}
			for _, is := range iss {
				sender.AddRows(key, is.Title, is.Link, is.CreatedAt)
				count++
			}
		}

		listID, err := s.stor.ListIDByName(nowLabel)
		if err != nil {
			return err
		}
		c, err := s.stor.CountIssuesByListID(listID)
		if err != nil {
			return err
		}

		startRow := c
		if startRow < 0 {
			log.Printf("warning: startRow %s is equal to %d\n", nowLabel, startRow)
			startRow = 2
		}

		sender.SetStartRange("A", startRangeRow+startRow)

		if err := s.svr.Update(sender, nowLabel); err != nil {
			return err
		}

		return s.stor.UpdateIncompleteIssues(processedOverdueIssues, closedIssues)
	}
}

func (s *Service) sheetChecker(listName string) error {
	exist, err := s.stor.ListExistByName(listName)
	if err != nil {
		return err
	}
	if !exist {
		if err := s.stor.InsertList(&storage.List{
			Name:      listName,
			CreatedAt: time.Now(),
		}); err != nil {
			return err
		}

		var err error
		if s.cfg.TmplSheetID == -1 {
			err = s.svr.CreateSheet(listName)
		} else {
			err = s.svr.CopySheet(s.cfg.TmplSheetID, listName)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// createdAfterTimeIssue gets the time from the first open issue of a specific cvs
func (s *Service) createdAfterTimeIssue(cvs string) (time.Time, error) {
	allCount, err := s.stor.CountIssuesByCVS(cvs)
	if err != nil {
		return time.Time{}, err
	}
	if allCount == 0 {
		return s.cfg.StartingTime, nil
	}
	openedCount, err := s.stor.CountOpenedIssuesByCVS(cvs)
	if err != nil {
		return time.Time{}, err
	}
	if openedCount == 0 {
		return s.cfg.StartingTime, nil
	}
	t, err := s.stor.FirstOpenedIssue(cvs)
	if err != nil {
		return time.Time{}, err
	}
	if t.After(s.cfg.StartingTime) {
		return t, nil
	}
	return s.cfg.StartingTime, nil
}
