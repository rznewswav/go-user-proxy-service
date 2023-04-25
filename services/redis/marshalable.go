package redis

type Marshalable interface {
	Marshal() (map[string]interface{}, error)
	Unmarshal(map[string]interface{}) error
}
