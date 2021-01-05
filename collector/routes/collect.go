package routes

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/pixel-collector/collector/middlewares"
	"github.com/spaghettifunk/pixel-collector/collector/models"
)

const (
	eventsTopic = "events.raw"
)

var (
	// used to fast unmarshal json strings
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// EventRequest is the request object for incoming messages
type EventRequest struct{}

// erPool is in charged of Pooling eventual requests in coming. This will help to reduce the alloc/s
// and effeciently improve the garbage collection operations. rr is short for event-request
var erPool = sync.Pool{
	New: func() interface{} { return new(EventRequest) },
}

// Collect receives the incoming request and forward the message to the bus after validation
func Collect(c echo.Context) error {
	// get a new object from the pool and then dispose it
	rr := erPool.Get().(interface{})
	defer erPool.Put(rr)

	// unmarshal GET query params into object
	raw := new(models.PayloadRaw)
	if err := c.Bind(raw); err != nil {
		log.Error().Err(err)
	}

	// enrich raw payload
	payload := models.Payload{
		ID:         uuid.New(),
		PayloadRaw: *raw,
	}

	// marshal payload for kafka
	p, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err)
	}

	// send it kafka
	pcc := c.(*middlewares.PixelContext)
	pcc.GetKafkaClient().Write(eventsTopic, p)

	return c.JSON(http.StatusOK, "OK")
}
