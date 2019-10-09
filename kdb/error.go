package kdb

import "github.com/pkg/errors"

var (
	ErrOutOfMemory         = errors.New("C01110 - OutOfMemory")
	// TODO REVIEW: Installation error missing
	ErrInternal            = errors.New("C01310 - Internal")
	ErrInterface           = errors.New("C01320 - Interface")
	ErrPluginMisbehavior   = errors.New("C01330 - PluginMisbehavior")
	ErrConflictingState    = errors.New("C02000 - ConflictingState")
	ErrValidationSyntactic = errors.New("C03100 - ValidationSyntactic")
	ErrValidationSemantic  = errors.New("C03200 - ValidationSemantic")
)

var (
	errCodeMap = map[string]error{
		"C01110": ErrOutOfMemory,
		"C01310": ErrInternal,
		"C01320": ErrInterface,
		"C01330": ErrPluginMisbehavior,
		"C02000": ErrConflictingState,
		"C03100": ErrValidationSyntactic,
		"C03200": ErrValidationSemantic,
	}
)
