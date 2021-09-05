package postgres

import (
	"compressor/internal/urlData"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

var _ urlData.URLDatas = &URLStorage{}

type URLStorage struct {
	statementStorage

	setURLStmt           *sql.Stmt
	setURLCompressedStmt *sql.Stmt
	getCompressedURLStmt *sql.Stmt
}

func NewURLStorage(db *DB) (*URLStorage, error) {
	s := &URLStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: setURLQuery, Dst: &s.setURLStmt},
		{Query: setURLCompressedQuery, Dst: &s.setURLCompressedStmt},
		{Query: getFullURLQuery, Dst: &s.getCompressedURLStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const setURLQuery = "INSERT INTO public.url (url) VALUES ($1) RETURNING id"

// SetURL ...
func (s *URLStorage) SetURL(url *urlData.URLData) error {
	if err := s.setURLStmt.QueryRow(&url.URL).Scan(&url.ID); err != nil {
		msg := fmt.Sprintf("can not exec query with url %v", url.URL)
		return errors.WithMessage(err, msg)
	}

	return nil
}

const setURLCompressedQuery = "UPDATE public.url SET url_compressed=$1 WHERE url=$2"

// SetURLCompressed ...
func (s *URLStorage) SetURLCompressed(url *urlData.URLData) error {
	if _, err := s.setURLCompressedStmt.Exec(&url.URLCompressed, &url.URL); err != nil {
		msg := fmt.Sprintf("can not exec query with url %v", url.URL)
		return errors.WithMessage(err, msg)
	}

	return nil
}

const getFullURLQuery = "SELECT url FROM public.url WHERE url_compressed=$1"

// GetFullURL ...
func (s *URLStorage) GetFullURL(url *urlData.URLData) error {
	row := s.getCompressedURLStmt.QueryRow(&url.URLCompressed)
	if err := scanURLData(row, url); err != nil {
		msg := fmt.Sprintf("can not exec query with url %v", url.URL)
		return errors.WithMessage(err, msg)
	}

	return nil
}

func scanURLData(scanner sqlScanner, url *urlData.URLData) error {
	return scanner.Scan(&url.URL)
}
