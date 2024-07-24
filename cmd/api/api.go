package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/faldeus0092/go-ecom/services/cart"
	"github.com/faldeus0092/go-ecom/services/order"
	"github.com/faldeus0092/go-ecom/services/product"
	"github.com/faldeus0092/go-ecom/services/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter) //register the user routes by passing the mux router

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	// run server, db not yet used
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}