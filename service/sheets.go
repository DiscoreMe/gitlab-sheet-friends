package service

import (
	"github.com/DiscoreMe/gitlab-sheets-friends/sheets"
)

// UpdateToken updates google api token
func (s *Service) UpdateToken() error {
	return s.svr.UpdateToken()
}

func (s *Service) UpdateColumns() error {
	sender := &sheets.Sender{}
	sender.SetStartRange("E", 2)

	var names []interface{}
	for _, member := range s.cfg.Members {
		names = append(names, member.Name)
	}
	sender.AddRows(names...)

	lists, err := s.stor.Lists()
	if err != nil {
		return err
	}
	for _, list := range lists {
		if err := s.svr.Update(sender, list.Name); err != nil {
			return err
		}
	}

	return nil
}
