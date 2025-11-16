package model

import (
	"time"
)

type Category int32

const (
	CategoryUnknown  Category = 0
	CategoryEngine   Category = 1
	CategoryFuel     Category = 2
	CategoryPorthole Category = 3
	CategoryWing     Category = 4
)

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width, omitempty"`
	Height float64 `bson:"height, omitempty"`
	Weight float64 `bson:"weight, omitempty"`
}

type Manufacturer struct {
	Name    string `bson:"name, omitempty"`
	Country string `bson:"country, omitempty"`
	Website string `bson:"website, omitempty"`
}

type Part struct {
	Uuid          string         `bson:"uuid"`
	Name          string         `bson:"name"`
	Description   string         `bson:"description"`
	Price         float64        `bson:"price"`
	StockQuantity int64          `bson:"stock_quantity"`
	Category      Category       `bson:"category"`
	Dimensions    *Dimensions    `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer  `bson:"manufacturer,omitempty"`
	Tags          []string       `bson:"tags,omitempty"`
	Metadata      map[string]any `bson:"metadata,omitempty"`
	CreatedAt     *time.Time     `bson:"created_at"`
	UpdatedAt     *time.Time     `bson:"updated_at,omitempty"`
}

type PartsFilter struct {
	Uuids                 []string   `bson:"uuids,omitempty"`
	Names                 []string   `bson:"names,omitempty"`
	Categories            []Category `bson:"categories,omitempty"`
	ManufacturerCountries []string   `bson:"manafacturer_countries,omitempty"`
	Tags                  []string   `bson:"tags, omitempty"`
}
