package storage

import (
	"database/sql/driver"
)

type Role string

const (
	Admin  Role = "admin"
	Viewer Role = "viewer"
)

func (r *Role) Scan(value interface{}) error {
	*r = Role(value.(string))
	return nil
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

func (r Role) String() string {
	return string(r)
}
