package storage

import (
	"database/sql"
	"time"
)

type Issue struct {
	ID        int       `db:"id"`
	Git       string    `db:"git"`
	ProjectID int       `db:"project_id"`
	IssueID   int       `db:"issue_id"`
	CreatedAt time.Time `db:"created_at"`
	ListID    int       `db:"list_id"`
	IsClosed  bool      `db:"is_closed"`
}

func (s *Storage) InsertIssue(issue *Issue) error {
	_, err := s.db.Exec(`INSERT INTO issues (git, project_id, issue_id, created_at, list_id, is_closed) VALUES ($1, $2, $3, $4, $5, $6);`, issue.Git, issue.ProjectID, issue.IssueID, issue.CreatedAt, issue.ListID, issue.IsClosed)
	return err
}
func (s *Storage) insertIssue(tx *sql.Tx, issue *Issue) error {
	_, err := tx.Exec(`INSERT INTO issues (git, project_id, issue_id, created_at, list_id, is_closed) VALUES ($1, $2, $3, $4, $5, $6);`, issue.Git, issue.ProjectID, issue.IssueID, issue.CreatedAt, issue.ListID, issue.IsClosed)
	return err
}

func (s *Storage) IssueByGitProjectIDIssueIDListID(git string, projectID, issueID, listID int) (Issue, error) {
	var issue Issue
	return issue, s.db.Get(&issue, `SELECT * FROM issues WHERE git = $1 AND project_id = $2 AND issue_id = $3 AND list_id = $4 LIMIT 1;`, git, projectID, issueID, listID)
}

func (s *Storage) LastCreatedTimeIssue(git string, listID int) (time.Time, error) {
	var t time.Time
	err := s.db.Get(&t, "SELECT created_at FROM issues WHERE git = $1 AND list_id = $2 ORDER BY created_at DESC LIMIT 1;", git, listID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return time.Time{}, err
	}
	return t, nil
}

func (s *Storage) CountIssuesByListID(id int) (int, error) {
	var c int
	err := s.db.QueryRow("SELECT count(*) FROM issues WHERE list_id = $1", id).Scan(&c)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return 0, err
	}
	return c, nil
}

func (s *Storage) CountIssuesByCVS(name string) (int, error) {
	var c int
	err := s.db.QueryRow("SELECT count(*) FROM issues WHERE git = $1", name).Scan(&c)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return 0, err
	}
	return c, nil
}

func (s *Storage) CountOpenedIssuesByCVS(name string) (int, error) {
	var c int
	err := s.db.QueryRow("SELECT count(*) FROM issues WHERE git = $1 AND is_closed = 0", name).Scan(&c)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return 0, err
	}
	return c, nil
}

func (s *Storage) FirstOpenedIssue(name string) (time.Time, error) {
	var t time.Time
	err := s.db.Get(&t, `select created_at from issues where git = $1 and created_at >= (select created_at from issues where git = $1 and is_closed = 0 order by created_at asc limit 1) order by created_at asc`, name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return time.Time{}, err
	}
	return t, nil
}

func (s *Storage) IssuesByIsClosed(isClosed bool) (issues []Issue, err error) {
	err = s.db.Select(&issues, "SELECT * FROM issues WHERE is_closed = $1", isClosed)
	return
}

func (s *Storage) UpdateIncompleteIssues(overdueIssues []Issue, closeIssues []Issue) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if closeIssues != nil {
		for _, issue := range closeIssues {
			if err := s.updateIssueIsClosedByGitProjectIDIssueIDCreatedAt(tx, issue); err != nil {
				return err
			}
		}
	}
	if overdueIssues != nil {
		for _, issue := range overdueIssues {
			if err := s.insertIssue(tx, &issue); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *Storage) updateIssueIsClosedByGitProjectIDIssueIDCreatedAt(tx *sql.Tx, issue Issue) error {
	_, err := tx.Exec("UPDATE issues SET is_closed = $1 WHERE git = $2 AND project_id = $3 AND issue_id = $4", issue.IsClosed, issue.Git, issue.ProjectID, issue.IssueID)
	return err
}
