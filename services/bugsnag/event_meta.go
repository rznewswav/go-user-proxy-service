package bugsnag

import driver "github.com/bugsnag/bugsnag-go/v2"

type EventMeta struct {
	Handler     string      // Handler type (kafka/internal emit)
	EventName   string      // Event name
	Payload     interface{} // Event payload
	TriggeredBy interface{} // Event meta
}

func (em *EventMeta) SetHandler(Handler string) *EventMeta {
	em.Handler = Handler
	return em
}

func (em *EventMeta) SetEventName(EventName string) *EventMeta {
	em.EventName = EventName
	return em
}

func (em *EventMeta) SetPayload(Payload any) *EventMeta {
	em.Payload = Payload
	return em
}

func (em *EventMeta) SetTriggeredBy(
	TriggeredBy any,
) *EventMeta {
	em.TriggeredBy = TriggeredBy
	return em
}

func (em *EventMeta) Decorate(
	event *driver.Event,
	config *driver.Configuration,
) {
	event.MetaData.AddStruct("event", em)
}
