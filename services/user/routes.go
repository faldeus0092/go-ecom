package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/faldeus0092/go-ecom/config"
	"github.com/faldeus0092/go-ecom/services/auth"
	"github.com/faldeus0092/go-ecom/types"
	"github.com/faldeus0092/go-ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore //we need userstore to interact with db
}

// make it same with Handler struct
func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request){
	// receive json payload, store it into payload
	var payload types.LoginUserPayload
	
	// parse
	if err := utils.ParseJSON(r, &payload); err != nil{ 
		utils.WriteError(w, http.StatusBadRequest, err) //doesn't print the error? it doesn't print the error because it's designed to write the error message to the HTTP response, not to the console
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil{
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	// check if user exists
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found, invalid email or password"))
		return
	}
	
	if !auth.ComparePasswords(u.Password, []byte(payload.Password)){
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request){
	// receive json payload, store it into payload
	var payload types.RegisterUserPayload
	
	// parse
	if err := utils.ParseJSON(r, &payload); err != nil{ 
		utils.WriteError(w, http.StatusBadRequest, err) //doesn't print the error? it doesn't print the error because it's designed to write the error message to the HTTP response, not to the console
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil{
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}
	
	// check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		// if there's no error, meaning email already exists
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}
	
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	// create new user otherwise
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Email: payload.Email,
		Password: hashedPassword,
	})
	
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	utils.WriteJSON(w, http.StatusCreated, nil)
	log.Println(err)
	//33:57
}