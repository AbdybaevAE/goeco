package codes

import (
	"net/http"

	"github.com/abdybaevae/goeco/pkg/opctx"
)

// Typed operation code(includes typed error that could appear)
type Code string

// Some known operation codes
const (
	Ok                 Code = "Ok"
	BadArguments       Code = "BadArguments"
	ServiceInternal    Code = "ServiceInternal"
	ServiceUnavailable Code = "ServiceUnavailable"
)

// mapping to http statuses
var codeToHttpStatus = map[Code]int{
	Ok:                 http.StatusOK,
	BadArguments:       http.StatusBadRequest,
	ServiceInternal:    http.StatusInternalServerError,
	ServiceUnavailable: http.StatusServiceUnavailable,
}

// default en locale code messages
var codeToMessage = map[Code]string{
	Ok:                 "Successfully proccessed.",
	BadArguments:       "Bad arguments provided.",
	ServiceInternal:    "Something wrong happened. Please, try again later.",
	ServiceUnavailable: "Something wrong happened. Please, try again later.",
}

// Every operation result. Includes code, message, id and data(for resource retrieve operations)
type OpRes struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Id      string      `json:"id"`
	Data    interface{} `json:"data,omitempty"`
}

// Default operation result factory(intercept all operation result creations)
type OpResFactory interface {
	// Create operation result from code, op.Op and human understandable message
	CreateWithMessage(op opctx.Op, code Code, message string) *OpRes
	// Create operation from code and op.Op
	Create(op opctx.Op, code Code) *OpRes
}

// default implementation of operation result factory
type opResFactoryImpl struct {
}

func (op *opResFactoryImpl) CreateWithMessage(opCtx opctx.Op, code Code, message string) *OpRes {
	return &OpRes{
		Code:    code,
		Message: message,
		Id:      opCtx.OpId(),
	}
}
func (op *opResFactoryImpl) Create(opCtx opctx.Op, code Code) *OpRes {
	message := codeToMessage[code]
	return op.CreateWithMessage(opCtx, code, message)
}
func New() OpResFactory {
	return &opResFactoryImpl{}
}
