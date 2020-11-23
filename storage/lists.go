package storage

import "time"

type List struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (s *Storage) InsertList(list *List) error {
	_, err := s.db.Exec(`INSERT INTO lists (name, created_at) VALUES ($1, $2)`, list.Name, list.CreatedAt)
	return err
}

func (s *Storage) Lists() ([]List, error) {
	var lists []List
	return lists, s.db.Select(&lists, "SELECT * FROM lists;")
}

func (s *Storage) ListIDByName(name string) (int, error) {
	var c int
	err := s.db.Get(&c, "SELECT id FROM lists WHERE name = $1", name)
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (s *Storage) ListExistByName(name string) (bool, error) {
	var c int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM lists WHERE name = $1", name).Scan(&c); err != nil {
		return false, err
	}
	if c > 0 {
		return true, nil
	}
	return false, nil
}
