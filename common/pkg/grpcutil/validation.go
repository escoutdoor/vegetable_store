package grpcutil

import (
	"errors"

	"buf.build/go/protovalidate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ProtoValidationError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := status.FromError(err); ok {
		return err
	}

	var validationErr *protovalidate.ValidationError
	ok := errors.As(err, &validationErr)
	if ok {
		st, stErr := status.New(codes.InvalidArgument, codes.InvalidArgument.String()).
			WithDetails(protovalidateErrToDetails(validationErr))
		if stErr == nil {
			return st.Err()
		}
	}

	return status.Error(codes.InvalidArgument, err.Error())
}

func protovalidateErrToDetails(validationErr *protovalidate.ValidationError) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: protovalidateViolationsToGoogleViolations(validationErr.Violations),
	}
}

func protovalidateViolationsToGoogleViolations(vs []*protovalidate.Violation) []*errdetails.BadRequest_FieldViolation {
	res := make([]*errdetails.BadRequest_FieldViolation, len(vs))
	for i, v := range vs {
		res[i] = &errdetails.BadRequest_FieldViolation{
			Field:       v.FieldDescriptor.JSONName(),
			Description: v.Proto.GetMessage(),
		}
	}
	return res
}
