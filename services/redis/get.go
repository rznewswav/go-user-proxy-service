package redis

import (
	"context"
	"encoding/json"
	"time"
)

func Get(unprefixedKey string, marshalable Marshalable) error {
	key := AddKeyPrefix(unprefixedKey)

	getContext, getCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer getCancel()

	jsonString, getError := Client.Get(getContext, key).Result()
	if getError != nil {
		return getError
	}

	var unmarshalled map[string]interface{}
	unmarshalError := json.Unmarshal([]byte(jsonString), &unmarshalled)
	if unmarshalError != nil {
		return unmarshalError
	}

	return marshalable.Unmarshal(unmarshalled)
}

func GetStruct[T any](unprefixedKey string, s *T) error {
	key := AddKeyPrefix(unprefixedKey)

	getContext, getCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer getCancel()

	jsonString, getError := Client.Get(getContext, key).Result()
	if getError != nil {
		return getError
	}

	unmarshalError := json.Unmarshal([]byte(jsonString), s)
	if unmarshalError != nil {
		return unmarshalError
	}
	return nil
}
