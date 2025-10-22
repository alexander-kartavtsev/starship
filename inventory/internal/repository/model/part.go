package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Uuid          string             `bson:"uuid"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	Price         float64            `bson:"price"`
	StockQuantity int64              `bson:"stockQuantity"`
	Category      Category           `bson:"category"`
	Dimensions    Dimensions         `bson:"dimensions"`
	Manufacturer  Manufacturer       `bson:"manufacturer"`
	Tags          []string           `bson:"tags"`
	CreatedAt     time.Time          `bson:"createdAt"`
	UpdatedAt     *time.Time         `bson:"updatedAt,omitempty"`
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	ManufacturerNames     []string
	Tags                  []string
}

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Category int32

const (
	CATEGORY_UNKNOWN Category = iota
	CATEGORY_ENGINE
	CATEGORY_FUEL
	CATEGORY_PORTHOLE
	CATEGORY_WING
)
