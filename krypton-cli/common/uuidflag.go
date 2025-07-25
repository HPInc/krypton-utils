// Package uuidflag wraps a UUID so it satisfies
// the Value and Getter interface of the Go flag package.
// Borrowed From: https://github.com/CSCfi/qvain-uuid on 12/28/2022
// MIT License
package common

import (
	uuid "github.com/google/uuid"
)

// Uuid wraps UUID.
type Uuid struct {
	uuid.UUID
	set bool
}

func NewUUID() *Uuid {
	return &Uuid{}
}

// Set function sets the UUID to the given string value or returns an error.
func (u *Uuid) Set(val string) (err error) {
	u.UUID, err = uuid.Parse(val)
	if err != nil {
		return err
	}
	u.set = true
	return nil
}

// String returns this type's default value for flag help.
// (If no default is set, the value is a zero-filled byte, so return empty string.)
func (u *Uuid) String() string {
	if u.set {
		return u.UUID.String()
	}
	return ""
}

// IsSet indicates whether the UUID value was set
// Returns false if there was an error setting value or if value was never set
func (u *Uuid) IsSet() bool {
	return u.set
}

// set default
func (u *Uuid) SetDefault() bool {
	return u.Set(uuid.NewString()) != nil
}
