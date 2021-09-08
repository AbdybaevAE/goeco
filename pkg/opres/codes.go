package opres

import (
	"net/http"
	"sync"
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

const (
	BadArgumentsMessage = "Bad arguments provided."
	GeneralErrorMessage = "Something wrong happened. Please, try again later."
	OkMessage           = "Successfully processed."
)

type Meta struct {
	HttpStatus int
	Message    string
}
type CodeMapper struct {
	mu       sync.Mutex
	mapStore map[Code]*Meta
}

var codeMapperInstance *CodeMapper
var codeMapperOnce sync.Once

func (cm *CodeMapper) Add(code Code, httpStatus int, message string) {
	cm.mu.Lock()
	if cm.mapStore == nil {
		cm.mapStore = make(map[Code]*Meta)
	}
	if _, ok := cm.mapStore[code]; ok {
		panic("already added given code " + string(code))
	}
	cm.mapStore[code] = &Meta{
		HttpStatus: httpStatus,
		Message:    message,
	}
	defer cm.mu.Unlock()
}
func (cm *CodeMapper) Get(code Code) *Meta {
	if val, ok := cm.mapStore[code]; !ok {
		panic("code is abset under code mapper")
	} else {
		return val
	}
}
func GetCodeMeta() *CodeMapper {
	codeMapperOnce.Do(func() {
		codeMapperInstance = &CodeMapper{}
		for _, metaItem := range defMeta {
			codeMapperInstance.Add(metaItem.Code, metaItem.HttpStatus, metaItem.Message)
		}
	})
	return codeMapperInstance
}
func init() {
	OpCodes = map[Code]*OpRes{}
	_ = GetCodeMeta()
	for _, meta := range defMeta {
		OpCodes[meta.Code] = &OpRes{
			Message: meta.Message,
			Code:    meta.Code,
		}
	}
}

var OpCodes map[Code]*OpRes

// default en locale code messages
var defMeta = []struct {
	Code       Code
	Message    string
	HttpStatus int
}{
	{Ok, OkMessage, http.StatusOK},
	{BadArguments, BadArgumentsMessage, http.StatusBadRequest},
	{ServiceInternal, GeneralErrorMessage, http.StatusInternalServerError},
	{ServiceUnavailable, GeneralErrorMessage, http.StatusServiceUnavailable},
}

var GeneralErrorOpCode = &OpRes{
	Code:    ServiceInternal,
	Message: GeneralErrorMessage,
}
var OkOpCode = &OpRes{
	Code:    Ok,
	Message: OkMessage,
}

// Every operation result. Includes code, message, id and data(for resource retrieve operations)
type OpRes struct {
	Code    Code        `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (o *OpRes) Error() string {
	return string(o.Code) + " " + o.Message
}

// Default operation result factory(intercept all operation result creations)
type OpResFactory interface {
	// Create operation result from code, op.Op and human understandable message
	CreateWithData(code Code, message string, data interface{}) *OpRes
	Create(code Code, message string) *OpRes
}

// default implementation of operation result factory
type opResFactoryImpl struct {
}

func (op *opResFactoryImpl) Create(code Code, message string) *OpRes {
	return op.CreateWithData(code, message, nil)
}
func NewOpRes() OpResFactory {
	return &opResFactoryImpl{}
}
func (op *opResFactoryImpl) CreateWithData(code Code, message string, data interface{}) *OpRes {
	return &OpRes{Code: code, Message: message, Data: data}
}
