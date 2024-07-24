package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/faldeus0092/go-ecom/types"
)

type Store struct {
	// dependency injection, so this depends on *sql.DB
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

func scanRowIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec("insert into products (name, description, image, price, quantity) values (?, ?, ?, ?, ?)", product.Name, product.Description, product.Image, product.Price, product.Quantity)
	if err != nil {
		return err
	}
	return nil
}

/*Accept an array of productIDs and returns an array of types.Product corresponding to the productIDs
 */
func (s *Store) GetProductsByIDs(productIDs []int) ([]types.Product, error) {
	// build query. appending ,? to already formatted ?%s
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("select * from products where id in (?%s)", placeholders)
	
	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	
	// need to return in the types.Product
	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	
	return products, nil

}

/* Update product based on types.Product received
 */
func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("update products set name =?, price=?, image=?, description=?, quantity=? where id=?", product.Name, product.Price, product.Image, product.Description, product.Quantity, product.ID)
	if err != nil {
		return err
	}
	return nil
}
