package bugsnag_structs

import (
	"encoding/json"
	"fmt"

	driver "github.com/bugsnag/bugsnag-go/v2"
)

type DatabaseQuery struct {
	Collection   string
	Operation    string
	Filter       any
	Sort         any
	Aggregations any
	BulkModels   any
	Options      any
	Document     any
}

func (dq *DatabaseQuery) SetCollection(
	Collection string,
) *DatabaseQuery {
	dq.Collection = Collection
	return dq
}

func (dq *DatabaseQuery) SetOperation(
	Operation string,
) *DatabaseQuery {
	dq.Operation = Operation
	return dq
}

func (dq *DatabaseQuery) SetFilter(Filter any) *DatabaseQuery {
	dq.Filter = Filter
	return dq
}

func (dq *DatabaseQuery) SetSort(Sort any) *DatabaseQuery {
	dq.Sort = Sort
	return dq
}

func (dq *DatabaseQuery) SetAggregations(
	Aggregations any,
) *DatabaseQuery {
	dq.Aggregations = Aggregations
	return dq
}

func (dq *DatabaseQuery) SetBulkModels(
	BulkModels any,
) *DatabaseQuery {
	dq.BulkModels = BulkModels
	return dq
}

func (dq *DatabaseQuery) SetOptions(
	Options any,
) *DatabaseQuery {
	dq.Options = Options
	return dq
}

func (dq *DatabaseQuery) SetDocument(
	Document any,
) *DatabaseQuery {
	dq.Document = Document
	return dq
}

func (dq *DatabaseQuery) Decorate(
	event *driver.Event,
	config *driver.Configuration,
) {
	event.MetaData.AddStruct("database", dq)
}

func (dq *DatabaseQuery) ToJSON() string {
	stringified, err := json.Marshal(dq)
	if err != nil {
		return fmt.Errorf("cannot stringify database query metadata: %w", err).
			Error()
	}
	return string(stringified)
}
