package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UUID          string             `bson:"uuid,omitempty"`
	Name          string             `bson:"name,omitempty"`
	Description   string             `bson:"description,omitempty"`
	Price         float64            `bson:"price,omitempty"`
	StockQuantity int64              `bson:"stock_quantity,omitempty"`
	Category      Category           `bson:"category,omitempty"`
	Dimensions    *Dimensions        `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer      `bson:"manufacturer,omitempty"`
	Tags          []string           `bson:"tags,omitempty"`
	Metadata      map[string]Value   `bson:"metadata,omitempty"`
	CreatedAt     *time.Time         `bson:"created_at,omitempty"`
	UpdatedAt     *time.Time         `bson:"updated_at,omitempty"`
}

type Dimensions struct {
	Length float64 `bson:"length,omitempty"`
	Width  float64 `bson:"width,omitempty"`
	Height float64 `bson:"height,omitempty"`
	Weight float64 `bson:"weight,omitempty"`
}

type Category string

const (
	ENGINE   Category = "ENGINE"
	FUEL     Category = "FUEL"
	PORTHOLE Category = "PORTHOLE"
	WING     Category = "WING"
)

type Manufacturer struct {
	Name    string `bson:"name,omitempty"`
	Country string `bson:"country,omitempty"`
	Website string `bson:"website,omitempty"`
}

type Value struct {
	StringValue  *string  `bson:"string_value,omitempty"`
	Int64Value   *int64   `bson:"int64_value,omitempty"`
	DoubleValue  *float64 `bson:"double_value,omitempty"`
	BooleanValue *bool    `bson:"boolean_value,omitempty"`
}

type PartsFilter struct {
	UUIDs                 []string   `bson:"uuids,omitempty"`
	Names                 []string   `bson:"names,omitempty"`
	Categories            []Category `bson:"categories,omitempty"`
	ManufacturerCountries []string   `bson:"manufacturer_countries,omitempty"`
	Tags                  []string   `bson:"tags,omitempty"`
}
