package models

import (
	"errors"
)

var (
	UserNotFoundErr     = errors.New("user not found")
	SegmentNotFoundErr  = errors.New("segment not found")
	MethodNotProvideErr = errors.New("method not provided")
	SegmentNameEmptyErr = errors.New("segment name is empty")
)
