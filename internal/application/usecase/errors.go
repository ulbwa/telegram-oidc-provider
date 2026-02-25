package usecase

import (
	"errors"
	"fmt"
)

var (
	// ErrUnexpected is returned for unexpected errors
	ErrUnexpected = errors.New("unexpected error occurred")

	// ErrInvalidInput is returned when input data is invalid or cannot be processed
	ErrInvalidInput = errors.New("invalid input data")
)

type GenericErr struct {
	Message string
}

func (e *GenericErr) Error() string { return e.Message }

type ObjectNotFoundErr struct {
	GenericErr
	Object string
	Id     any
}

func NewObjectNotFoundErr(object string, id any) error {
	err := new(ObjectNotFoundErr)
	err.Object = object
	if id == nil {
		err.Message = fmt.Sprintf("%s not found", object)
	} else {
		err.Message = fmt.Sprintf("%s with id '%v' not found", object, id)
		err.Id = id
	}
	return err
}

type ObjectInvalidErr struct {
	GenericErr
	Object string
	Field  string
	Reason *string
}

func NewObjectInvalidErr(object string, field string, reason *string) error {
	err := new(ObjectInvalidErr)
	err.Object = object
	err.Field = field
	if reason == nil {
		err.Message = fmt.Sprintf("%s has invalid field '%s'", object, field)
	} else {
		err.Message = fmt.Sprintf("%s has invalid field '%s': %s", object, field, *reason)
		err.Reason = reason
	}
	return err
}

type ConflictErr struct {
	GenericErr
	Object  string
	Feature *string
}

func NewConflictErr(object string, feature *string) error {
	err := new(ConflictErr)
	err.Object = object
	if feature == nil {
		err.Message = fmt.Sprintf("%s is already assigned", object)
	} else {
		err.Message = fmt.Sprintf("%s is already assigned to another %s", *feature, object)
		err.Feature = feature
	}
	return err
}

type BadGatewayErr struct {
	GenericErr
	Service string
}

func NewBadGatewayErr(service string) error {
	err := new(BadGatewayErr)
	err.Service = service
	err.Message = fmt.Sprintf("%s is unavailable", service)
	return err
}

type GatewayTimeoutErr struct {
	GenericErr
	Service string
}

func NewGatewayTimeoutErr(service string) error {
	err := new(GatewayTimeoutErr)
	err.Service = service
	err.Message = fmt.Sprintf("%s request timed out", service)
	return err
}
