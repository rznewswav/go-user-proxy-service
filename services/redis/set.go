package redis

import (
	"context"
	"encoding/json"
	"time"
)

type ExpiryConfig time.Duration

func Set(unprefixedKey string, marshalable Marshalable, opts ...interface{}) error {
	key := AddKeyPrefix(unprefixedKey)

	var expiration time.Duration = 0
	for _, opt := range opts {
		switch casted := opt.(type) {
		case ExpiryConfig:
			expiration = time.Duration(casted)
		}
	}
	marshalled, marshalError := marshalable.Marshal()
	if marshalError != nil {
		return marshalError
	}

	jsonMarshalled, jsonMarshalError := json.Marshal(marshalled)
	if jsonMarshalError != nil {
		return jsonMarshalError
	}
	jsonString := string(jsonMarshalled)

	setContext, setContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer setContextCancel()

	return Client.Set(setContext, key, jsonString, expiration).Err()
}

func SetStruct(unprefixedKey string, s interface{}, opts ...interface{}) error {
	key := AddKeyPrefix(unprefixedKey)

	var expiration time.Duration = 0
	for _, opt := range opts {
		switch casted := opt.(type) {
		case ExpiryConfig:
			expiration = time.Duration(casted)
		}
	}

	jsonMarshalled, jsonMarshalError := json.Marshal(s)
	if jsonMarshalError != nil {
		return jsonMarshalError
	}
	jsonString := string(jsonMarshalled)

	setContext, setContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer setContextCancel()

	return Client.Set(setContext, key, jsonString, expiration).Err()
}
