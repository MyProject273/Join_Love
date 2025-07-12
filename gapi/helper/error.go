package helper

import (
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func InvalidArrgumentError(violation []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violation}
	statusInvalid := status.New(codes.InvalidArgument, "invalid argument")

	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}

func ValidateAll(req interface{}) error {
	validator, ok := req.(interface {
		ValidateAll() error
	})
	if !ok {
		return status.Errorf(codes.InvalidArgument, "Request does not implement ValidateAll")
	}

	if err := validator.ValidateAll(); err != nil {
		var details []string

		switch v := err.(type) {
		case interface{ Unwrap() []error }:
			for _, subErr := range v.Unwrap() {
				if verr, ok := subErr.(interface {
					Error() string
				}); ok {
					details = append(details, verr.Error())
				}
			}
		default:
			details = append(details, v.Error())
		}
		return status.Errorf(codes.InvalidArgument, "Validation failed: %s", strings.Join(details, "; "))
	}
	return nil
}
