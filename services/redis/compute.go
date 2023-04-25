package redis

import "github.com/redis/go-redis/v9"

func GetOrCompute[T any](key string, computeFunction func() (T, error), saveOpts ...interface{}) (T, error) {
	var s T
	var getError error
	if marshalable, castable := any(s).(Marshalable); castable {
		getError = Get(key, marshalable)
	} else {
		getError = GetStruct(key, &s)
	}

	if getError == nil {
		return s, nil
	}

	if getError != redis.Nil {
		return s, getError
	}

	// value is nil in redis, compute and save
	computed, computeError := computeFunction()
	if computeError != nil {
		return s, computeError
	}

	if marshalable, castable := any(computed).(Marshalable); castable {
		Set(key, marshalable, saveOpts...)
	} else {
		SetStruct(key, computed, saveOpts...)
	}
	return computed, nil
}
