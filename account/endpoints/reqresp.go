package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	// CreateUserRequest createUserRequest
	CreateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// CreateUserResponse CreateUserResponse
	CreateUserResponse struct {
		OK string `json:"ok"`
	}

	// GetUserRequest GetUserRequest
	GetUserRequest struct {
		ID string `json:"id"`
	}

	// GetUserResponse GetUserResponse
	GetUserResponse struct {
		Email string `jsong:"email"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func decodeGetUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetUserRequest

	vars := mux.Vars(r)
	fmt.Printf("decode createUserReq: %s\n", vars["id"])
	req = GetUserRequest{
		ID: vars["id"],
	}
	return req, nil
}
