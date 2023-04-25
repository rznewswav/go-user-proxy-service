package controllers

type Headers struct {
	H *map[string]string
}

func (s *Headers) AssignMapIfNil() (m map[string]string) {
	if s.H != nil {
		m = *s.H
		return
	} else {
		m = make(map[string]string)
		s.H = &m
		return
	}
}

func (s *Headers) Set(key string, value string) {
	s.AssignMapIfNil()

	unreffedMap := *s.H
	unreffedMap[key] = value
}

func (s *Headers) SetterFunc() func(string, string) {
	return func(key, value string) {
		s.Set(key, value)
	}
}

func (s *Headers) Get(key string) (value string) {
	s.AssignMapIfNil()

	unreffedMap := *s.H
	ok := false
	if value, ok = unreffedMap[key]; !ok {
		value = ""
	}

	return
}

func (s *Headers) ForEach(iterator func(string, string)) {
	s.AssignMapIfNil()

	unreffedMap := *s.H

	for key, value := range unreffedMap {
		iterator(key, value)
	}
}
