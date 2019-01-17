package templates

type Visitor interface {
	VisitString(value string) (interface{}, error)
	VisitMapStringOfString(value map[string]string) (interface{}, error)
	VisitMapStringOfInterface(value map[string]interface{}) (interface{}, error)
	VisitMapInterfaceOfInterface(value map[interface{}]interface{}) (interface{}, error)
	VisitInterface(value interface{}) (interface{}, error)
	VisitSliceOfString(value []string) (interface{}, error)
	VisitSliceOfInterface(value []interface{}) (interface{}, error)
}
