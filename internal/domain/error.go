package domain

import "errors"

var (
	ErrTaskNotFound       = errors.New("task doesn't exist")
	ErrAssignmentNotFound = errors.New("assignment doesn't exist")
)
