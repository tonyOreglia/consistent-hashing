package controller

import "errors"

var (
	ErrNodeAlreadyExists = errors.New("node already exists")
	ErrNodeDoesNotExist  = errors.New("node does not exist")
)
