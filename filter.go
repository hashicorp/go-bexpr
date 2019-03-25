package filter

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-filter/bexpr"
)

type Filter struct {
	// The underlying boolean expression evaluator
	evaluator *bexpr.Evaluator
}

type CreateArgs struct {
	Expression string
	Config     bexpr.EvaluatorConfig
	dataType   interface{}
}

// Creates a filter to operate on the given data type.
// The data type is the type of the elements that will be filtered and not the top level container type.
// For example, if you want to filter a []Foo then the data type to pass here is Foo.
func Create(expression string, config *bexpr.EvaluatorConfig, dataType interface{}) (*Filter, error) {
	// TODO (mkeeler) - figure out how to allow getting rid of the top level map[*] or slice/array and getting
	// just the elem. I think this will require allowing the bexpr code to take reflect.Type instead of
	// unconditionally grabbing the type via reflect.TypeOf
	exp, err := bexpr.CreateForType(expression, config, dataType)
	if err != nil {
		return nil, fmt.Errorf("Failed to create boolean expression evaluator: %v", err)
	}

	return &Filter{
		evaluator: exp,
	}, nil
}

func (f *Filter) Execute(data interface{}) (interface{}, error) {
	rvalue := reflect.ValueOf(data)
	rtype := rvalue.Type()

	switch rvalue.Kind() {
	case reflect.Array:
		// For arrays we return slices instead of fixed sized arrays
		rtype = reflect.SliceOf(rtype.Elem())
		fallthrough
	case reflect.Slice:
		newSlice := reflect.MakeSlice(rtype, 0, rvalue.Len())

		for i := 0; i < rvalue.Len(); i++ {
			item := rvalue.Index(i)
			if !item.CanInterface() {
				return nil, fmt.Errorf("Slice/Array value can not be used")
			}
			result, err := f.evaluator.Evaluate(item.Interface())
			if err != nil {
				return nil, err
			}

			if result {
				newSlice = reflect.Append(newSlice, item)
			}
		}

		return newSlice.Interface(), nil
	case reflect.Map:
		newMap := reflect.MakeMap(rtype)

		// TODO (mkeeler) - Update to use a MapRange iterator once Go 1.12 is usable
		// for all of our products
		for _, mapKey := range rvalue.MapKeys() {
			item := rvalue.MapIndex(mapKey)

			if !item.CanInterface() {
				return nil, fmt.Errorf("Map value cannot be used")
			}

			result, err := f.evaluator.Evaluate(item.Interface())
			if err != nil {
				return nil, err
			}

			if result {
				newMap.SetMapIndex(mapKey, item)
			}
		}

		return newMap.Interface(), nil
	default:
		return nil, fmt.Errorf("Only slices, arrays and maps are filterable")
	}
}
