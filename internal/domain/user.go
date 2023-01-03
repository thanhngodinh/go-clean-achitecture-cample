package domain

import (
	"context"
	"net/http"
	"time"

	"github.com/core-go/search"
)

type User struct {
	Id          string     `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"-" avro:"id" validate:"required,max=40" match:"equal"`
	Username    string     `json:"username" gorm:"column:username" bson:"username" dynamodbav:"username" firestore:"username" avro:"username" validate:"required,username,max=100" match:"prefix"`
	Email       string     `json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" avro:"email" validate:"email,max=100" match:"prefix"`
	Phone       string     `json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" avro:"phone" validate:"required,phone,max=18"`
	DateOfBirth *time.Time `json:"dateOfBirth" gorm:"column:date_of_birth" bson:"dateOfBirth" dynamodbav:"dateOfBirth" firestore:"dateOfBirth" avro:"dateOfBirth"`
}

type UserFilter struct {
	*search.Filter
	Id          string            `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"id" avro:"id" validate:"required,max=40" match:"equal"`
	Username    string            `json:"username" gorm:"column:username" bson:"username" dynamodbav:"username" firestore:"username" avro:"username" validate:"required,username,max=100" match:"prefix" q:"prefix"`
	Email       string            `json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" avro:"email" validate:"email,max=100" match:"prefix" q:"prefix"`
	Phone       string            `json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" avro:"phone" validate:"required,phone,max=18" q:"true"`
	DateOfBirth *search.TimeRange `json:"dateOfBirth" gorm:"column:date_of_birth" bson:"dateOfBirth" dynamodbav:"dateOfBirth" firestore:"dateOfBirth" avro:"dateOfBirth"`
}

type UserHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type UserUsecase interface {
	Load(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (int64, error)
	Update(ctx context.Context, user *User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

type UserRepository interface {
	Load(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (int64, error)
	Update(ctx context.Context, user *User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
