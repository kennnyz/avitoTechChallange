package models

import (
	"errors"
)

var (
	UserNotFoundErr     = errors.New("user not found")
	MethodNotProvideErr = errors.New("method not provided")
	SegmentNameEmptyErr = errors.New("segment name is empty")
	NoDataErr           = errors.New("no data")
)

func SegmentNotFoundErr(segmentName string) error {
	return errors.New("segment " + segmentName + " not found")
}
