package cart

import (
	"fmt"
	"net/http"

	"github.com/faldeus0092/go-ecom/services/auth"
	"github.com/faldeus0092/go-ecom/types"
	"github.com/faldeus0092/go-ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.OrderStore
	productStore types.ProductStore // for checking product stock
	userStore types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) (*Handler){
	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router)  {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	// will receive an array of CartItem, contains ID and qty
	// need to convert these into array of types.Product
	// then usesr both cart.Items and []types.Product to make an order
	var cart types.CartCheckoutPayload

	// user ID is obtainable through JWT token
	userID := auth.GetUserIDFromContext(r.Context())
	
	// parse
	if err := utils.ParseJSON(r, &cart); err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate
	if err := utils.Validate.Struct(cart); err != nil{
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	// turn cart item IDs into array of IDs
	productIDs, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// query into array of types.Product
	products, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	
	// create new order and create every order items
	orderID, totalPrice, err := h.createOrder(products, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"order_id": orderID,
		"total_price": totalPrice,
	})
}