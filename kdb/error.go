package kdb

import "github.com/pkg/errors"

var (
	ErrOutOfMemory         = errors.New("C01110")
	ErrInternal            = errors.New("C01310")
	ErrInterface           = errors.New("C01320")
	ErrPluginMisbehavior   = errors.New("C01330")
	ErrConflictingState    = errors.New("C02000")
	ErrValidationSyntactic = errors.New("C03100")
	ErrValidationSemantic  = errors.New("C03200")
)

var (
	ErrCodeMap = map[string]error{
		"C01110": ErrOutOfMemory,
		"C01310": ErrInternal,
		"C01320": ErrInterface,
		"C01330": ErrPluginMisbehavior,
		"C02000": ErrConflictingState,
		"C03100": ErrValidationSyntactic,
		"C03200": ErrValidationSemantic,
	}
)
