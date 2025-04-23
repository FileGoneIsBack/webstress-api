package users

import (
  "api/core/database"
  "errors"
)

type Instance struct {
	conn *database.Instance
  }
  
  func New(db *database.Instance) *Instance {
	return &Instance{conn: db}
}
type Query interface{ Scan(...any) error }
var (
	Container = New(database.Container)
	ErrDuplicateUser = errors.New("duplicate user")
	ErrUserNotFound  = errors.New("user couldn't be found in the database")
	ErrInvalidInput  = errors.New("invalid ticket data")
)