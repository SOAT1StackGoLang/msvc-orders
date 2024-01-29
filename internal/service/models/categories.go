package models

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

type CategoryList struct {
	Categories    []*Category
	Limit, Offset int
	Total         int64
}
