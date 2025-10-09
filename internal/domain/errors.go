package domain

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input data")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrConflict     = errors.New("resource conflict")
	ErrInternal     = errors.New("internal server error")

	ErrRoomUnavailable  = errors.New("room not available for the selected time slot")
	ErrTimeRangeInvalid = errors.New("invalid start or end time for booking")
)
