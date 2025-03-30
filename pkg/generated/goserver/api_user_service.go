// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Diary - OpenAPI 3.0
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 * Contact: ilya.korolev@outlook.com
 */

package goserver

import (
	"context"
	"errors"
	"net/http"
)

// UserAPIService is an interface that defines the logic for the UserAPIServicer
type UserAPIService interface {
	// GetUser - return user object
	GetUser(ctx context.Context) (ImplResponse, error)
}

// UserAPIService is a service that implements the logic for the UserAPIServicer
// This service should implement the business logic for every endpoint for the UserAPI API.
// Include any external packages or services that will be required by this service.
type UserAPIServiceImpl struct {
}

// NewUserAPIService creates a default api service
func NewUserAPIService() UserAPIService {
	return &UserAPIServiceImpl{}
}

// GetUser - return user object
func (s *UserAPIServiceImpl) GetUser(ctx context.Context) (ImplResponse, error) {
	// TODO - update GetUser with the required logic for this service method.
	// Add api_user_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, User{}) or use other options such as http.Ok ...
	// return Response(200, User{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetUser method not implemented")
}
