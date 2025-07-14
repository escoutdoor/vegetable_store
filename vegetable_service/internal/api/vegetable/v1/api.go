package v1

import (
	"fmt"

	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/service"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	minVegetableNameLength = 3
	maxVegetableNameLength = 25

	nameFieldName            = "name"
	weightFieldName          = "weight"
	priceFieldName           = "price"
	discountedPriceFieldName = "discounted_price"
)

type Implementation struct {
	vegetableService service.VegetableService
	vegetablev1.UnimplementedVegetableServiceServer
}

func NewImplementation(vegetableService service.VegetableService) *Implementation {
	return &Implementation{
		vegetableService: vegetableService,
	}
}

func validateBatchUpdateVegetablesRequest(req *vegetablev1.BatchUpdateVegetablesRequest) error {
	for _, r := range req.Requests {
		if err := validateUpdateVegetableRequest(r); err != nil {
			return err
		}
	}

	return nil
}

func validateUpdateVegetableRequest(req *vegetablev1.UpdateVegetableRequest) error {
	if req.UpdateMask == nil {
		return fmt.Errorf("update_mask: unspecified")
	}
	if len(req.UpdateMask.Paths) == 0 {
		return fmt.Errorf("update_mask: paths: unspecified")
	}

	em := make(map[string]error, len(req.UpdateMask.Paths))
	for _, p := range req.UpdateMask.Paths {
		switch p {
		case nameFieldName:
			if err := validateVegetableName(req.Update.Name); err != nil {
				em[nameFieldName] = err
			}
		case weightFieldName:
			if err := validatePositiveFloat(req.Update.Weight); err != nil {
				em[weightFieldName] = err
			}
		case priceFieldName:
			if err := validatePositiveFloat(req.Update.Price); err != nil {
				em[priceFieldName] = err
			}
		case discountedPriceFieldName:
			if err := validatePositiveFloat(req.Update.DiscountedPrice); err != nil {
				em[discountedPriceFieldName] = err
			}
		}
	}
	if len(em) > 0 {
		vs := make([]*errdetails.BadRequest_FieldViolation, 0, len(em))
		for f, e := range em {
			vs = append(vs, &errdetails.BadRequest_FieldViolation{
				Field:       f,
				Description: e.Error(),
			})
		}

		st := status.New(codes.InvalidArgument, codes.InvalidArgument.String())

		badReq := &errdetails.BadRequest{FieldViolations: vs}

		detailedStatus, err := st.WithDetails(badReq)
		if err != nil {
			return err
		}

		return detailedStatus.Err()
	}

	return nil
}

func validateVegetableName(name string) error {
	if len(name) < minVegetableNameLength {
		return fmt.Errorf("value length must be at least %d characters", minVegetableNameLength)
	}

	if len(name) > maxVegetableNameLength {
		return fmt.Errorf("value length must be at most %d characters", maxVegetableNameLength)
	}

	return nil
}

func validatePositiveFloat(value float32) error {
	if value <= 0 {
		return fmt.Errorf("value must be greater than 0")
	}

	return nil
}
