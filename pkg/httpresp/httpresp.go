package httpresp

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/abdybaevae/goeco/pkg/codes"
	"github.com/abdybaevae/goeco/pkg/msg"
	opres "github.com/abdybaevae/goeco/pkg/op/res"
	"github.com/sirupsen/logrus"
)

type RsFactory interface {
	// Default response with data
	OkData(rw http.ResponseWriter, data interface{})
	// Success repsonse status
	Ok(rw http.ResponseWriter)
	// Send response from OpRes variable(includes code)
	Send(rw http.ResponseWriter, opRes *opres.Res)
	// Error response handler
	Error(rw http.ResponseWriter, err error)
	// Response from code
	Code(rw http.ResponseWriter, code codes.Code)
	// Response from code and message
	CodeMessage(rw http.ResponseWriter, code codes.Code, message string)
}
type rsFactoryImpl struct {
	logger *logrus.Logger
}

var rsFactoryInstance RsFactory
var rsFactoryOnce sync.Once
var cachedRes = opres.GetCacheRes()

func Get() RsFactory {
	rsFactoryOnce.Do(func() {
		rsFactoryInstance = &rsFactoryImpl{
			logger: logrus.New(),
		}
	})
	return rsFactoryInstance
}
func (r *rsFactoryImpl) Send(rw http.ResponseWriter, opRes *opres.Res) {
	metaData := cachedRes.GetMeta(opRes.Code)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(metaData.HttpStatus)
	json.NewEncoder(rw).Encode(opRes)
}

var opResFactory = opres.NewOpRes()

func (r *rsFactoryImpl) Error(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "application/json")
	if opErr, ok := err.(*opres.Res); ok {
		r.Send(rw, opErr)
	} else {
		r.logger.Warnf("unhandler error %v", err)
		r.Send(rw, cachedRes.Get(codes.Ok))
	}
}
func (r *rsFactoryImpl) Ok(rw http.ResponseWriter) {
	r.Send(rw, cachedRes.Get(codes.Ok))
}
func (r *rsFactoryImpl) OkData(rw http.ResponseWriter, data interface{}) {
	r.Send(rw, opResFactory.CreateWithData(codes.Ok, msg.OkMsgCode, data))
}
func (r *rsFactoryImpl) Code(rw http.ResponseWriter, code codes.Code) {
	r.Send(rw, cachedRes.Get(code))
}
func (r *rsFactoryImpl) CodeMessage(rw http.ResponseWriter, code codes.Code, message string) {
	opRes := opResFactory.Create(code, message)
	r.Send(rw, opRes)
}
