package structs

type DefaultedMap[K comparable, V any] struct {
	H                *map[K]V
	defaultWhenEmpty V
}

type StringDefaultedMap = DefaultedMap[string, string]

//goland:noinspection GoUnusedExportedFunction
func NewDefaultedMap[K comparable, V any](defaultWhenEmpty V) DefaultedMap[K, V] {
	return DefaultedMap[K, V]{defaultWhenEmpty: defaultWhenEmpty}
}

//goland:noinspection GoUnusedExportedFunction
func NewStringMap() DefaultedMap[string, string] {
	return DefaultedMap[string, string]{defaultWhenEmpty: ""}
}

func (s *DefaultedMap[K, V]) AssignMapIfNil() (m map[K]V) {
	if s.H != nil {
		m = *s.H
		return
	} else {
		m = make(map[K]V)
		s.H = &m
		return
	}
}

func (s *DefaultedMap[K, V]) Set(key K, value V) {
	s.AssignMapIfNil()

	unreffedMap := *s.H
	unreffedMap[key] = value
}

func (s *DefaultedMap[K, V]) SetterFunc() func(K, V) {
	return func(key K, value V) {
		s.Set(key, value)
	}
}

func (s *DefaultedMap[K, V]) Get(key K) (value V) {
	s.AssignMapIfNil()

	unreffedMap := *s.H
	ok := false
	if value, ok = unreffedMap[key]; !ok {
		value = s.defaultWhenEmpty
	}

	return
}

func (s *DefaultedMap[K, V]) Has(key K) (exists bool) {
	s.AssignMapIfNil()

	unreffedMap := *s.H
	_, exists = unreffedMap[key]
	return
}

func (s *DefaultedMap[K, V]) ForEach(iterator func(K, V)) {
	s.AssignMapIfNil()

	unreffedMap := *s.H

	for key, value := range unreffedMap {
		iterator(key, value)
	}
}
