// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type AuthPayload struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type AuthResponse struct {
	AuthToken *AuthToken `json:"authToken"`
	User      *User      `json:"user"`
}

type AuthToken struct {
	AccessToken string    `json:"accessToken"`
	ExpiredAt   time.Time `json:"expiredAt"`
}

type CreateProductInput struct {
	Category    ProductCategory `json:"category"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
}

type CreateUserInput struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Password string   `json:"password"`
	City     string   `json:"city"`
	Mode     UserMode `json:"mode"`
}

type LoginInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Mutation struct {
}

type Query struct {
}

type Subscription struct {
}

type ProductCategory string

const (
	ProductCategoryHome    ProductCategory = "HOME"
	ProductCategoryClothes ProductCategory = "CLOTHES"
)

var AllProductCategory = []ProductCategory{
	ProductCategoryHome,
	ProductCategoryClothes,
}

func (e ProductCategory) IsValid() bool {
	switch e {
	case ProductCategoryHome, ProductCategoryClothes:
		return true
	}
	return false
}

func (e ProductCategory) String() string {
	return string(e)
}

func (e *ProductCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProductCategory(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProductCategory", str)
	}
	return nil
}

func (e ProductCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProductStatus string

const (
	ProductStatusCreated ProductStatus = "CREATED"
	ProductStatusSold    ProductStatus = "SOLD"
)

var AllProductStatus = []ProductStatus{
	ProductStatusCreated,
	ProductStatusSold,
}

func (e ProductStatus) IsValid() bool {
	switch e {
	case ProductStatusCreated, ProductStatusSold:
		return true
	}
	return false
}

func (e ProductStatus) String() string {
	return string(e)
}

func (e *ProductStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProductStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProductStatus", str)
	}
	return nil
}

func (e ProductStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UserMode string

const (
	UserModeClient UserMode = "CLIENT"
	UserModeAdmin  UserMode = "ADMIN"
)

var AllUserMode = []UserMode{
	UserModeClient,
	UserModeAdmin,
}

func (e UserMode) IsValid() bool {
	switch e {
	case UserModeClient, UserModeAdmin:
		return true
	}
	return false
}

func (e UserMode) String() string {
	return string(e)
}

func (e *UserMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserMode", str)
	}
	return nil
}

func (e UserMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
