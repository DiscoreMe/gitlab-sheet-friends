package sheets

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
	"net/http"
)

const (
	credentialsPath = "credentials.json"
	tokenPath       = "token.json"
)

const scopeSheets = "https://www.googleapis.com/auth/spreadsheets"

type ErrTokenNotFoundOrOutdated struct {
	Err error
}

func (e ErrTokenNotFoundOrOutdated) Error() string {
	return fmt.Sprintf("token not found or is outdated: %s", e.Err)
}

type Sheets struct {
	id     string
	client *http.Client
	svr    *sheets.Service
}

func NewSheets(id string) (*Sheets, error) {
	s := &Sheets{id: id}

	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Sheets) init() error {
	srvConfig, err := parseConfig(credentialsPath)
	if err != nil {
		return err
	}

	client, err := initClient(tokenPath, srvConfig)
	if err != nil {
		return err
	}
	s.client = client

	svr, err := sheets.New(s.client)
	if err != nil {
		return err
	}
	s.svr = svr

	return nil
}

func initClient(tokenPath string, config *oauth2.Config) (*http.Client, error) {
	tokFile := tokenPath
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, ErrTokenNotFoundOrOutdated{
			Err: err,
		}
	}
	client := config.Client(context.Background(), tok)
	return client, nil
}

func (s *Sheets) Update(sender *Sender, list string) error {
	_, err := s.svr.Spreadsheets.Values.Update(s.id, list+"!"+sender.StartRange(), &sheets.ValueRange{
		Values: sender.Rows(),
	}).ValueInputOption("USER_ENTERED").Context(context.TODO()).Do()
	return err
}

func (s *Sheets) CreateSheet(name string) error {
	_, err := s.svr.Spreadsheets.BatchUpdate(s.id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title:  name,
					Hidden: false,
				},
			}},
		},
	}).Context(context.TODO()).Do()
	return err
}

func (s *Sheets) CopySheet(srcID int64, name string) error {
	_, err := s.svr.Spreadsheets.BatchUpdate(s.id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{DuplicateSheet: &sheets.DuplicateSheetRequest{
				SourceSheetId: srcID,
				NewSheetName:  name,
			}},
		},
	}).Context(context.TODO()).Do()
	return err
}
