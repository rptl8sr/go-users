package ownErrors

import (
	"fmt"
)

// ErrNotFound is used to indicate that a requested record could not be found in the database.
var ErrNotFound = fmt.Errorf("record not found")

var ErrUserAlreadyExists = fmt.Errorf("user already exists")
