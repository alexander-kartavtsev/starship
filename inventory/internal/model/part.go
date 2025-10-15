package model

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
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
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Category int

const (
	CATEGORY_UNKNOWN Category = iota
	CATEGORY_ENGINE
	CATEGORY_FUEL
	CATEGORY_PORTHOLE
	CATEGORY_WING
)
