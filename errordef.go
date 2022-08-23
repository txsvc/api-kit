package apikit

import "errors"

const (
	MsgStatus = "status: %d"
)

var (
	// ErrInvalidConfiguration indicates that parameters used to configure the service were invalid
	ErrInvalidConfiguration = errors.New("invalid configuration")
	// ErrInvalidPassPhrase indicates that the pass phrase is too short
	ErrInvalidPassPhrase = errors.New("invalid pass phrase")
	// ErrAlreadyInitialized indicates that client is already registered
	ErrAlreadyInitialized = errors.New("already initialized")
	// ErrInitializingConfiguration indicates that the client could not be initialized
	ErrInitializingConfiguration = errors.New("error initializing")

	// ErrNotImplemented indicates that a function is not yet implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrInternalError indicates everything else
	ErrInternalError = errors.New("internal error")
	// ErrApiError indicates an error in an API call
	ErrApiError = errors.New("api error")
	// ErrInvalidRoute indicates that the route and/or its parameters are not valid
	ErrInvalidRoute = errors.New("invalid route")
	// ErrInvalidResourceName indicates that the resource name is invalid
	ErrInvalidResourceName = errors.New("invalid resource name")
	// ErrMissingResourceName indicates that a resource type is missing
	ErrMissingResourceName = errors.New("missing resource type")
	// ErrResourceNotFound indicates that the resource does not exist
	ErrResourceNotFound = errors.New("resource does not exist")
	// ErrResourceExists indicates that the resource does not exist
	ErrResourceExists = errors.New("resource already exists")
	// ErrInvalidParameters indicates that parameters used in an API call are not valid
	ErrInvalidParameters = errors.New("invalid parameters")
	// ErrInvalidNumArguments indicates that the number of arguments in an API call is not valid
	ErrInvalidNumArguments = errors.New("invalid arguments")
)
