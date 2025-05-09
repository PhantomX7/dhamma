package utility

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/PhantomX7/dhamma/utility/logger"
)

// ContextValues holds values extracted from the context.
type ContextValues struct {
	DomainID *uint64
	UserID   uint64
	IsRoot   bool
}

// NewContextWithValues creates a new context with the provided ContextValues.
func NewContextWithValues(ctx context.Context, values ContextValues) context.Context {
	return context.WithValue(ctx, "values", values)
}

// ValuesFromContext retrieves ContextValues from the given context.
// It returns an error if the values are not found or are of an unexpected type.
func ValuesFromContext(ctx context.Context) (ContextValues, error) {
	values, ok := ctx.Value("values").(ContextValues)
	if !ok {
		logger.Get().Error("failed to retrieve values from context")
		return ContextValues{}, &errors.AppError{
			Message: "forbidden",
			Status:  http.StatusBadRequest,
		}
	}
	return values, nil
}

// CheckDomainContext validates if an operation on a specific entity is allowed for the domain specified in the context.
// It compares the domainID from the context with the provided entityDomainID.
// entityName is used to create a specific error message (e.g., "follower", "product").
// actionVerb is used to create a more specific error message (e.g., "create", "update", "get").
// Returns the ContextValues and an error if the domain check fails or if context values cannot be retrieved.
func CheckDomainContext(ctx context.Context, entityDomainID uint64, entityName string, actionVerb string) (ContextValues, error) {
	contextValues, err := ValuesFromContext(ctx)
	if err != nil {
		// This indicates an issue with retrieving context values, which is an internal problem.
		return ContextValues{}, err
	}

	// If DomainID is present in the context, enforce the check.
	if contextValues.DomainID != nil {
		if entityDomainID != *contextValues.DomainID {
			return contextValues, &errors.AppError{
				Message: fmt.Sprintf("you cannot %s %s for another domain", actionVerb, entityName),
				Status:  http.StatusBadRequest,
			}
		}
	}
	// If contextValues.DomainID is nil, it implies no domain restriction from the context, so the operation is permitted from a domain perspective.
	// Or, if DomainIDs match, the operation is permitted.
	return contextValues, nil
}
