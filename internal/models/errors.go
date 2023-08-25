package models

import (
	"errors"
)

var (
	UserNotFound        = errors.New("user not found")
	SegmentNotFound     = errors.New("segment not found")
	MethodNotProvideErr = errors.New("method not provided")
)
