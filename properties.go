package propertygraph2go

type WithProperties interface {
	GetProperty(key string) (value interface{}, err error)
	SetProperty(key string, value interface{}) (err error)
	Properties()(properties map[string]interface{}, err error)
}

func HasProperty(w WithProperties, key string) (has bool) {
	_, err := w.GetProperty(key)
	has = err == nil
	return
}

func HasPropertyValue(w WithProperties, key string, value interface{}) (has bool) {
	v, err := w.GetProperty(key)
	has = err == nil
	if !has {
		return
	}
	has = v == value
	return
}
