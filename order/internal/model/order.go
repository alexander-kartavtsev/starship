package model

type OrderInfo struct {
	UserUuid  string
	PartUuids []string
}

func (i *OrderInfo) GetUserUuid() string {
	return i.UserUuid
}

func (i *OrderInfo) GetPartUuids() []string {
	return i.PartUuids
}

type OrderUpdateInfo struct {
	TransactionUuid *string
	PartUuids       *[]string
	PaymentMethod   *PaymentMethod
	Status          *OrderStatus
}

type Order struct {
	OrderUuid       string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUuid string
	PaymentMethod   PaymentMethod
	Status          OrderStatus
}

type OrderCreateRes struct {
	OrderUuid  string
	TotalPrice float64
}

type PaymentMethod string

const (
	Unknown       PaymentMethod = "UNKNOWN"
	card          PaymentMethod = "CARD"
	sbp           PaymentMethod = "SBP"
	creditCard    PaymentMethod = "CREDIT_CARD"
	investorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	PendingPayment OrderStatus = "PENDING_PAYMENT"
	Paid           OrderStatus = "PAID"
	Cancelled      OrderStatus = "CANCELLED"
)

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	ManufacturerNames     []string
	Tags                  []string
}

type Part struct {
	Uuid          string
	Price         float64
	StockQuantity int64
}

type Category int

const (
	CATEGORY_UNKNOWN Category = iota
	CATEGORY_ENGINE
	CATEGORY_FUEL
	CATEGORY_PORTHOLE
	CATEGORY_WING
)
