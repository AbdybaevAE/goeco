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

type meta struct {
	HttpStatus int
}

var cachedCodesInstanceOnce sync.Once
var cm *cacheCodesImpl

func GetCacheRes() CachedRes {
	cachedCodesInstanceOnce.Do(func() {
		cm = &cacheCodesImpl{}
		cm.MetaMap = make(map[codes.Code]*meta)
		cm.ResMap = make(map[codes.Code]*Res)
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
			cm.ResMap[item.Code] = &Res{
				Message: item.Message,
				Code:    item.Code,
			}
			cm.MetaMap[item.Code] = &meta{
				HttpStatus: item.HttpStatus,
			}
		}
	})
	return cm
}

type CachedRes interface {
	Add(code codes.Code, httpStatus int, message string)
	Get(code codes.Code) *Res
	GetMeta(code codes.Code) *meta
}
type cacheCodesImpl struct {
	once    sync.Once
	mu      sync.Mutex
	MetaMap map[codes.Code]*meta
	ResMap  map[codes.Code]*Res
}

func (cm *cacheCodesImpl) Add(code codes.Code, httpStatus int, message string) {
	cm.mu.Lock()
	if _, ok := cm.MetaMap[code]; ok {
		panic("already added given code " + string(code))
	}
	cm.ResMap[code] = &Res{
		Message: message,
		Code:    code,
	}
	cm.MetaMap[code] = &meta{
		HttpStatus: httpStatus,
	}
	defer cm.mu.Unlock()
}
func (cm *cacheCodesImpl) Get(code codes.Code) *Res {
	if val, ok := cm.ResMap[code]; !ok {
		panic("no code " + string(code) + " found")
	} else {
		return val
	}
}
func (cm *cacheCodesImpl) GetMeta(code codes.Code) *meta {
	return cm.MetaMap[code]
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
