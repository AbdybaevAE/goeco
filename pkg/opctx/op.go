package opctx

import "github.com/sirupsen/logrus"

// Operation interface that hold current operation data
// requestId - current request id, unique per each incoming request
// operationId - current operation id. If request need's to be processed by multiple microservices(which has suitable probabality) then
// each of microservice requests has same operation id.
// It's very usefull for logging and tracing
// Operation has logger with operationId and requestId included.
type Op interface {
	ReqId() string
	OpId() string
	Log() *logrus.Entry
}

const (
	opKey  = "opId"
	reqKey = "reqId"
)

// default implementation
type opImpl struct {
	RequestId   string
	OperationId string
	logger      *logrus.Entry
}

func (o *opImpl) ReqId() string {
	return o.RequestId
}
func (o *opImpl) OpId() string {
	return o.OperationId
}
func (o *opImpl) Log() *logrus.Entry {
	return o.logger
}

// constructor
func New(requestId, operationId string) Op {
	if len(requestId) == 0 {
		panic("cannot create operation from empty requestId")
	}
	if len(operationId) == 0 {
		panic("cannot create operation from empty operationId")
	}
	logger := logrus.WithFields(logrus.Fields{
		reqKey: requestId,
		opKey:  opKey,
	})
	return &opImpl{
		RequestId:   requestId,
		OperationId: operationId,
		logger:      logger,
	}
}
