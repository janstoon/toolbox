package bricks

import "errors"

var (
	ErrResolvable = errors.New("temporary issue")

	// ----------------------
	//     Customer-side
	// ----------------------

	ErrCustomerSide       = errors.New("customer-side error")
	ErrLogical            = errors.Join(ErrCustomerSide, errors.New("logical issue"))
	ErrInvalidInput       = errors.Join(ErrCustomerSide, errors.New("invalid input"))
	ErrInvalidCredentials = errors.Join(ErrInvalidInput, errors.New("invalid credentials"))
	ErrPermissionDenied   = errors.Join(ErrCustomerSide, errors.New("permission denied"))
	ErrNotFound           = errors.Join(ErrCustomerSide, errors.New("not found"))
	ErrReachedEnd         = errors.Join(ErrCustomerSide, errors.New("reached end"))

	// ----------------------
	//     Supplier-side
	// ----------------------

	ErrSupplierSide   = errors.New("supplier-side error")
	ErrInfrastructure = errors.Join(ErrSupplierSide, ErrResolvable, errors.New("infrastructure"))
	ErrNotImplemented = errors.Join(ErrSupplierSide, errors.New("not implemented"))
)

type MessageEnvelope struct {
	Id string

	Retried int
	Note    []byte
}
