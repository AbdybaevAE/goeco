package httpresp

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/abdybaevae/goeco/pkg/opres"
	"github.com/sirupsen/logrus"
)

type RsFactory interface {
	// Default response with data
	OkData(rw http.ResponseWriter, data interface{})
	// Success repsonse status
	Ok(rw http.ResponseWriter)
	// Send response from OpRes variable(includes code)
	Send(rw http.ResponseWriter, opRes *opres.OpRes)
	// Error response handler
	Error(rw http.ResponseWriter, err error)
	// Response from code
	Code(rw http.ResponseWriter, code opres.Code)
	// Response from code and message
	CodeMessage(rw http.ResponseWriter, code opres.Code, message string)
}
type rsFactoryImpl struct {
	logger *logrus.Logger
}

var meta = opres.GetCodeMeta()
var rsFactoryInstance RsFactory
var rsFactoryOnce sync.Once

func Get() RsFactory {
	rsFactoryOnce.Do(func() {
		rsFactoryInstance = &rsFactoryImpl{
			logger: logrus.New(),
		}
	})
	return rsFactoryInstance
}
func (r *rsFactoryImpl) Send(rw http.ResponseWriter, opRes *opres.OpRes) {
	metaData := meta.Get(opRes.Code)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(metaData.HttpStatus)
	json.NewEncoder(rw).Encode(opRes)
}

var opResFactory = opres.NewOpRes()

func (r *rsFactoryImpl) Error(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "application/json")
	if opErr, ok := err.(*opres.OpRes); ok {
		r.Send(rw, opErr)
	} else {
		r.logger.Warnf("unhandler error %v", err)
		r.Send(rw, opres.OpCodes[opres.Ok])
	}
}
func (r *rsFactoryImpl) Ok(rw http.ResponseWriter) {
	r.Send(rw, opres.OpCodes[opres.Ok])
}
func (r *rsFactoryImpl) OkData(rw http.ResponseWriter, data interface{}) {
	r.Send(rw, opResFactory.CreateWithData(opres.Ok, opres.OkOpCode.Message, data))
}
func (r *rsFactoryImpl) Code(rw http.ResponseWriter, code opres.Code) {
	r.Send(rw, opres.OpCodes[code])
}
func (r *rsFactoryImpl) CodeMessage(rw http.ResponseWriter, code opres.Code, message string) {
	opRes := opResFactory.Create(code, message)
	r.Send(rw, opRes)
}
