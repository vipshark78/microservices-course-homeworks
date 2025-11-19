package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode представляет коды ошибок.
type ErrorCode int64

const (
	NotFoundErrCode ErrorCode = iota
	InvalidArgumentErrCode
)

// businessError представляет бизнес ошибку.
type businessError struct {
	code ErrorCode
	err  error
}

func (b *businessError) Error() string {
	if b.err != nil {
		return b.err.Error()
	}
	return "unknown business error"
}

func (b *businessError) Unwrap() error {
	return b.err
}

func (b *businessError) Code() ErrorCode {
	return b.code
}

// NewNotFoundError создает новую бизнес ошибку с кодом NotFoundErrCode.
func NewNotFoundError(err error) *businessError {
	return &businessError{
		code: NotFoundErrCode,
		err:  err,
	}
}

// NewInvalidArgumentError создает новую бизнес ошибку с кодом InvalidArgumentErrCode.
func NewInvalidArgumentError(err error) *businessError {
	return &businessError{
		code: InvalidArgumentErrCode,
		err:  err,
	}
}

// GetBusinessError возвращает бизнес ошибку из переданной ошибки или nil, если ошибка не является бизнес ошибкой.
func GetBusinessError(err error) *businessError {
	var businessErr *businessError
	if errors.As(err, &businessErr) {
		return businessErr
	}
	return nil
}

// ErrorCodeToGRPCCode получает код ошибки и возвращает соответствующий код gRPC
func ErrorCodeToGRPCCode(code ErrorCode) codes.Code {
	switch code {
	case NotFoundErrCode:
		return codes.NotFound
	case InvalidArgumentErrCode:
		return codes.InvalidArgument
	default:
		return codes.Unknown
	}
}

// BusinessErrorToGRPCStatus конвертирует бизнес ошибку в статус gRPC
func BusinessErrorToGRPCStatus(err *businessError) *status.Status {
	grpcCode := ErrorCodeToGRPCCode(err.Code())
	return status.New(grpcCode, err.Error())
}
