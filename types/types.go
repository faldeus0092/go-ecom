package types

import "time"

type UserStore interface{
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user User) error
}

// for interacting with DB, make it same with table in DB
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

// for register json payload
type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=130"`
}

// for login json payload
type LoginUserPayload struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type ProductStore interface{
	GetProducts()([]Product, error)
	GetProductsByIDs(products []int) ([]Product, error)
	CreateProduct(product Product) error
	UpdateProduct(product Product) error
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int      `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

// for create product payload
type CreateProductPayload struct{
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Quantity    int      `json:"quantity" validate:"required"`
}

type OrderStore interface{
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}

type Order struct{
	ID int `json:"id"`
	UserID int `json:"userID"`
	Total float64 `json:"total"`
	Status string `json:"status"`
	Address string `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct{
	ID int `json:"id"`
	OrderID int `json:"orderID"`
	ProductID int `json:"productID"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type CartItem struct{
	ProductID int `json:"productID"`
	Quantity int `json:"quantity"`
}

type CartCheckoutPayload struct{
	Items []CartItem `json:"items" validate:"required"`
}