package res

import (
	"net/http"
	"sync"

	"github.com/abdybaevae/goeco/pkg/codes"
	"github.com/abdybaevae/goeco/pkg/msg"
)

// Every operation result. Includes code, message, id and data(for resource retrieve operations)
type Res struct {
	Code    codes.Code  `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (o *Res) Error() string {
	return string(o.Code) + " " + o.Message
}

type allCodes struct {
	once    sync.Once
	mu      sync.Mutex
	metaMap map[codes.Code]*meta
	resMap  map[codes.Code]*Res
}
type meta struct {
	HttpStatus int
}

var GenRes *allCodes

func (cm *allCodes) Init() {
	cm.once.Do(func() {
		cm.metaMap = make(map[codes.Code]*meta)
		cm.resMap = make(map[codes.Code]*Res)
		for _, item := range []struct {
			Code       codes.Code
			Message    string
			HttpStatus int
		}{
			{codes.Ok, msg.OkMsgCode, http.StatusOK},
			{codes.BadArguments, msg.BadArgumentsMsgCode, http.StatusBadRequest},
			{codes.ServiceInternal, msg.GeneralErrorMsgCode, http.StatusInternalServerError},
			{codes.ServiceUnavailable, msg.GeneralErrorMsgCode, http.StatusServiceUnavailable},
		} {
			cm.resMap[item.Code] = &Res{
				Message: item.Message,
				Code:    item.Code,
			}
			cm.metaMap[item.Code] = &meta{
				HttpStatus: item.HttpStatus,
			}
		}
	})
}
func (cm *allCodes) Add(code codes.Code, httpStatus int, message string) {
	cm.mu.Lock()
	if _, ok := cm.metaMap[code]; ok {
		panic("already added given code " + string(code))
	}
	cm.resMap[code] = &Res{
		Message: message,
		Code:    code,
	}
	cm.metaMap[code] = &meta{
		HttpStatus: httpStatus,
	}
	defer cm.mu.Unlock()
}
func (cm *allCodes) Get(code codes.Code) *Res {
	if val, ok := cm.resMap[code]; !ok {
		panic("no code " + string(code) + " found")
	} else {
		return val
	}
}
func (cm *allCodes) GetMeta(code codes.Code) *meta {
	return cm.metaMap[code]
}
func init() {
	GenRes = &allCodes{}
	GenRes.Init()
}

// Default operation result factory(intercept all operation result creations)
type OpResFactory interface {
	// Create operation result from code, op.Op and human understandable message
	CreateWithData(code codes.Code, message string, data interface{}) *Res
	Create(code codes.Code, message string) *Res
}

// default implementation of operation result factory
type opResFactoryImpl struct {
}

func (op *opResFactoryImpl) Create(code codes.Code, message string) *Res {
	return op.CreateWithData(code, message, nil)
}
func NewOpRes() OpResFactory {
	return &opResFactoryImpl{}
}
func (op *opResFactoryImpl) CreateWithData(code codes.Code, message string, data interface{}) *Res {
	return &Res{Code: code, Message: message, Data: data}
}
