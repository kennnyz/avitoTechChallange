package models

import (
	"errors"
	"fmt"
)

var (
	UserNotFoundErr     = errors.New("user not found")
	MethodNotProvideErr = errors.New("method not provided")
	SegmentNameEmptyErr = errors.New("segment name is empty")
	NoDataErr           = errors.New("no data")
	SegmentNotFoundErr  = errors.New("segment not found")
)

func NewSegmentNotFoundErr(segmentName string) error {
	return fmt.Errorf("%w, segment name %s", SegmentNotFoundErr, segmentName)
}
