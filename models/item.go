package models

import "github.com/gocql/gocql"

type Item struct {
    ItemID   gocql.UUID `json:"item_id"`
    Name     string     `json:"name"`
    Category string     `json:"category"`
}

type ItemPopularity struct {
    ItemID           gocql.UUID `json:"item_id"`
    PopularityCounter int64      `json:"popularity_counter"`
}
