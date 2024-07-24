package order

import (
	"database/sql"
	"fmt"

	"github.com/faldeus0092/go-ecom/types"
)

type Store struct {
	// dependency injection, so this depends on *sql.DB
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	res, err := s.db.Exec("insert into orders (userId, total, status, address) values (?, ?, ?, ?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("insert into order_items (orderId, productId, quantity, price) values (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}

func (s *Store) UpdateOrder(order types.Order) error {
	_, err := s.db.Exec("update orders set userId = ?, total = ?, status = ?, address = ? where id = ?", order.UserID, order.Total, order.Status, order.Address, order.ID)
	return err
}

func (s *Store) GetOrdersByUserID(userID int) ([]types.Order, error) {
	rows, err := s.db.Query("SELECT * FROM orders WHERE userId = ?", userID)
	if err != nil {
		return nil, err
	}
	orders := make([]types.Order, 0)
	for rows.Next(){
		order, err := scanRowIntoOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (s *Store) GetOrderByID(orderID int) (*types.Order, error){
	rows, err := s.db.Query("SELECT * FROM orders WHERE id = ?", orderID)
	if err != nil {
		return nil, err
	}
	order := new(types.Order)
	for rows.Next(){
		order, err = scanRowIntoOrder(rows)
		if err != nil {
			return nil, err
		}
	}

	if order.ID == 0{
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

func scanRowIntoOrder(rows *sql.Rows) (*types.Order, error) {
	order := new(types.Order)
	err := rows.Scan(&order.ID,
		&order.UserID,
		&order.Total,
		&order.Status,
		&order.Address,
		&order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}
