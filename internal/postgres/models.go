// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package postgres

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoleType string

const (
	RoleTypeAdmin RoleType = "admin"
	RoleTypeStaff RoleType = "staff"
	RoleTypeUser  RoleType = "user"
)

func (e *RoleType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RoleType(s)
	case string:
		*e = RoleType(s)
	default:
		return fmt.Errorf("unsupported scan type for RoleType: %T", src)
	}
	return nil
}

type NullRoleType struct {
	RoleType RoleType
	Valid    bool // Valid is true if RoleType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRoleType) Scan(value interface{}) error {
	if value == nil {
		ns.RoleType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RoleType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRoleType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RoleType), nil
}

type PointsOfInterest struct {
	ID        uuid.UUID
	Name      string
	Key       string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type User struct {
	ID          uuid.UUID
	Email       string
	Password    string
	MfaSecret   pgtype.Text
	MfaEnabled  pgtype.Bool
	MfaVerified pgtype.Bool
	Role        RoleType
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}
