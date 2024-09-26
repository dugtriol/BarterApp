// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/dugtriol/BarterApp/graph/scalar"
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
	AccessToken string          `json:"accessToken"`
	ExpiredAt   scalar.DateTime `json:"expiredAt"`
}

type ChangeStatusInput struct {
	ID     string            `json:"id"`
	Status TransactionStatus `json:"status"`
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

type Favorites struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
}

type LoginInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Mutation struct {
}

type Query struct {
}

type Transaction struct {
	ID             string              `json:"id"`
	Owner          string              `json:"owner"`
	Buyer          string              `json:"buyer"`
	ProductIDOwner string              `json:"product_id_owner"`
	ProductIDBuyer string              `json:"product_id_buyer"`
	CreatedAt      scalar.DateTime     `json:"created_at"`
	Shipping       TransactionShipping `json:"shipping"`
	Address        string              `json:"address"`
	Status         TransactionStatus   `json:"status"`
}

type TransactionCreateInput struct {
	Owner          string              `json:"owner"`
	ProductIDOwner string              `json:"product_id_owner"`
	ProductIDBuyer string              `json:"product_id_buyer"`
	Shipping       TransactionShipping `json:"shipping"`
	Address        string              `json:"address"`
}

type ProductCategory string

const (
	ProductCategoryHome     ProductCategory = "HOME"
	ProductCategoryClothes  ProductCategory = "CLOTHES"
	ProductCategoryChildren ProductCategory = "CHILDREN"
	ProductCategorySport    ProductCategory = "SPORT"
	ProductCategoryOther    ProductCategory = "OTHER"
)

var AllProductCategory = []ProductCategory{
	ProductCategoryHome,
	ProductCategoryClothes,
	ProductCategoryChildren,
	ProductCategorySport,
	ProductCategoryOther,
}

func (e ProductCategory) IsValid() bool {
	switch e {
	case ProductCategoryHome, ProductCategoryClothes, ProductCategoryChildren, ProductCategorySport, ProductCategoryOther:
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
	ProductStatusAvailable  ProductStatus = "AVAILABLE"
	ProductStatusExchanging ProductStatus = "EXCHANGING"
	ProductStatusExchanged  ProductStatus = "EXCHANGED"
)

var AllProductStatus = []ProductStatus{
	ProductStatusAvailable,
	ProductStatusExchanging,
	ProductStatusExchanged,
}

func (e ProductStatus) IsValid() bool {
	switch e {
	case ProductStatusAvailable, ProductStatusExchanging, ProductStatusExchanged:
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

type TransactionShipping string

const (
	TransactionShippingMeetup  TransactionShipping = "MEETUP"
	TransactionShippingMail    TransactionShipping = "MAIL"
	TransactionShippingCourier TransactionShipping = "COURIER"
)

var AllTransactionShipping = []TransactionShipping{
	TransactionShippingMeetup,
	TransactionShippingMail,
	TransactionShippingCourier,
}

func (e TransactionShipping) IsValid() bool {
	switch e {
	case TransactionShippingMeetup, TransactionShippingMail, TransactionShippingCourier:
		return true
	}
	return false
}

func (e TransactionShipping) String() string {
	return string(e)
}

func (e *TransactionShipping) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TransactionShipping(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TransactionShipping", str)
	}
	return nil
}

func (e TransactionShipping) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TransactionStatus string

const (
	TransactionStatusCreated  TransactionStatus = "CREATED"
	TransactionStatusOngoing  TransactionStatus = "ONGOING"
	TransactionStatusDone     TransactionStatus = "DONE"
	TransactionStatusDeclined TransactionStatus = "DECLINED"
)

var AllTransactionStatus = []TransactionStatus{
	TransactionStatusCreated,
	TransactionStatusOngoing,
	TransactionStatusDone,
	TransactionStatusDeclined,
}

func (e TransactionStatus) IsValid() bool {
	switch e {
	case TransactionStatusCreated, TransactionStatusOngoing, TransactionStatusDone, TransactionStatusDeclined:
		return true
	}
	return false
}

func (e TransactionStatus) String() string {
	return string(e)
}

func (e *TransactionStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TransactionStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TransactionStatus", str)
	}
	return nil
}

func (e TransactionStatus) MarshalGQL(w io.Writer) {
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
