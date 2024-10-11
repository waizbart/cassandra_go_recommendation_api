package models

import "github.com/gocql/gocql"

type Item struct {
    ItemID   gocql.UUID `json:"item_id"`
    Name     string     `json:"name"`
    Category string     `json:"category"`
}
