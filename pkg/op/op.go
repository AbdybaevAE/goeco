package op

import "github.com/sirupsen/logrus"

// Operation interface that hold current operation data
// requestId - current request id, unique per each incoming request
// operationId - current operation id. If request need's to be processed by multiple microservices(which has suitable probabality) then
// each of microservice requests has same operation id.
// It's very usefull for logging and tracing
// Operation has logger with operationId and requestId included.
type Operation interface {
	Request() string
	Operation() string
	Log() string
}

// default implementation
type operationImpl struct {
	RequestId   string
	OperationId string
	// logger
}

func (o *operationImpl) Request() string {
	return o.RequestId
}
func (o *operationImpl) Operation() string {
	return o.OperationId
}
func (o *operationImpl) Log() {

}

// constructor
func New(requestId, operationId string) Operation {
	if len(requestId) == 0 {
		panic("cannot create operation from empty requestId")
	}
	if len(operationId) == 0 {
		panic("cannot create operation from empty operationId")
	}
	log := logrus.WithFields(logrus.Fields{
		"asd": "21321",
	})
	return &operationImpl{
		RequestId:   requestId,
		OperationId: operationId,
		// logger: logger,
	}
}
