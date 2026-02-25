package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ulbwa/telegram-oidc-provider/api/generated"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
	"github.com/ulbwa/telegram-oidc-provider/pkg/utils"
)

func handleError(err error) (int, *generated.ErrorResponse, error) {
	var objNotFoundErr *usecase.ObjectNotFoundErr
	if errors.As(err, &objNotFoundErr) {
		var resp generated.ErrorResponse
		resp.Message = objNotFoundErr.Error()
		var details generated.ObjectNotFoundDetails
		details.Type = generated.ObjectNotFound
		if objNotFoundErr.Id != nil {
			details.Id = utils.Ptr(fmt.Sprintf("%v", objNotFoundErr.Id))
		}
		details.Object = objNotFoundErr.Object
		resp.Details = &generated.ErrorResponse_Details{}
		resp.Details.FromObjectNotFoundDetails(details)
		return http.StatusNotFound, &resp, nil
	}

	var objInvalidErr *usecase.ObjectInvalidErr
	if errors.As(err, &objInvalidErr) {
		var resp generated.ErrorResponse
		resp.Message = objInvalidErr.Error()
		var details generated.ObjectInvalidDetails
		details.Type = generated.ObjectInvalid
		details.Object = objInvalidErr.Object
		details.Field = objInvalidErr.Field
		details.Reason = objInvalidErr.Reason
		resp.Details = &generated.ErrorResponse_Details{}
		resp.Details.FromObjectInvalidDetails(details)
		return http.StatusBadRequest, &resp, nil
	}

	var conflictErr *usecase.ConflictErr
	if errors.As(err, &conflictErr) {
		var resp generated.ErrorResponse
		resp.Message = conflictErr.Error()
		var details generated.ConflictDetails
		details.Type = generated.Conflict
		details.Object = conflictErr.Object
		details.Feature = conflictErr.Feature
		resp.Details = &generated.ErrorResponse_Details{}
		resp.Details.FromConflictDetails(details)
		return http.StatusConflict, &resp, nil
	}

	if errors.Is(err, usecase.ErrInvalidInput) {
		var resp generated.ErrorResponse
		resp.Message = err.Error()
		return http.StatusBadRequest, &resp, nil
	}

	if errors.Is(err, usecase.ErrUnexpected) {
		var resp generated.ErrorResponse
		resp.Message = err.Error()
		return http.StatusInternalServerError, &resp, nil
	}

	return 0, nil, err
}
