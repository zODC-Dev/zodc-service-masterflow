package services

import "database/sql"

type HistoryService struct {
	db *sql.DB
}

func NewHistoryService(db *sql.DB) *HistoryService {
	return &HistoryService{
		db: db,
	}
}
