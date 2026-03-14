package counter

import errors "github.com/AntonioMartinezLopez/enginsight/pkg"

const (
	ErrTypeDomain = "domain"
)

const (
	ErrCodeInternal = iota + 5000000
)

// Here, sentinel errors can be added if needed

func NewInternalError(msg string, underlyingErr error) *errors.Error {
	return &errors.Error{
		Message: msg,
		Code:    ErrCodeInternal,
		Type:    ErrTypeDomain,
		Err:     underlyingErr,
	}
}
