package model

import "errors"


var ErrPartNotFound = errors.New("part not found")

var ErrPartUUIDIsEmpty = errors.New("part uuid is empty")

var ErrPartsNotFound = errors.New("parts with this filter not found")

var ErrUUIDIsNotValid = errors.New("part uuid is not uuid format")