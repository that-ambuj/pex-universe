package model

import (
	"reflect"
	"time"
)

type AutoInc struct {
	Id        uint64     `db:"id" json:"id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

type PaginationDto struct {
	Page  int
	Limit int
}

func (p *PaginationDto) Skip() int {
	return (p.Page - 1) * p.Limit
}

type PaginatedResponse struct {
	Data        interface{}
	CurrentPage int
	TotalPages  int
}

func IsNonZero(val *reflect.Value) bool {
	return val.CanUint() && !val.IsZero()
}

func IsNotEmptyString(val *reflect.Value) bool {
	return val.Kind() == reflect.String && val.String() != ""
}
