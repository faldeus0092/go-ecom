package cart

import (
	"fmt"
	"github.com/faldeus0092/go-ecom/types"
)

/*	Check if cartItems have enough stock for each product
*	returns an array of product IDs
 */
func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	products := make([]int, len(items))
	for i, v := range items {
		// check if quantity ordered is valid
		if v.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %v", v.ProductID)
		}
		products[i] = v.ProductID
	}
	return products, nil
}

/* Create order based on array of 
*	returns order id, total price, and error
*/
func (h *Handler) createOrder(products []types.Product, cart types.CartCheckoutPayload, userID int) (int, float64, error){
	// for convenience
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}
	
	// check if all products in stock
	if err := checkIfCartIsInStock(cart.Items, productMap); err != nil{
		return 0, 0, err
	}
	// calculate the total price
	totalPrice := calculateTotalPrice(cart.Items, productMap)
	
	// reduce quantity of products in our db
	for _, item := range cart.Items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.productStore.UpdateProduct(product)
	}
	
	// create the order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID: userID,
		Total: totalPrice,
		Status: "pending", //todo
		Address: cart.Address, //todo
	})
	if err != nil {
		return 0, 0, err
	}

	// create order cart.Items
	for _, item := range cart.Items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID: orderID,
			ProductID: item.ProductID,
			Quantity: item.Quantity,
			Price: productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	// cartItems => contains product id and bought quantity
	// products => contains product data stored in DB
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product with id %d not available, please refresh cart", item.ProductID)
		}
		if item.Quantity > int(product.Quantity){
			return fmt.Errorf("insufficient stock for product %s", product.Name)
		}

	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	var total float64 = 0.00
	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price*float64(item.Quantity)
	}
	return total
}