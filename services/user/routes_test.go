package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/faldeus0092/go-ecom/types"
	"github.com/gorilla/mux"
)

func TestUserServiceHandler(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T){
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName: "asdf",
			Email: "asd",
			Password: "asdfgasdfasdf",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		// should be fail
		log.Println("\nhere!")
		log.Println(rr.Code)
		if rr.Code != http.StatusBadRequest{
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName: "asdf",
			Email: "asd@gmail.com",
			Password: "asdfgasdfasdf",
		}
		marshalled, _ := json.Marshal(payload)
		
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		// should be fail
		log.Println("\nhere!")
		log.Println(rr.Code)
		if rr.Code != http.StatusCreated{
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockUserStore struct{}

// implement mockUserStore the same as UserStore in types.go
func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error){
	return nil, nil
}

func (m *mockUserStore) CreateUser(user types.User) error{
	return nil
}
