package propertygraph


type GraphItem interface {
	Id()         string
	Properties() interface{}
}
