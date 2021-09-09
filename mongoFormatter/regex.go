package mongoFormatter

import(
	"strings"
	"fmt"
	"reflect"
	"errors"
)

var funcMapper = map[string]interface{}{
	"startsWith": startsWith,
	"endsWith": endsWith,
	"contains": contains,
	"equals": equals,
}

func startsWith(val string, regex string) bool{
	return strings.HasPrefix(val, regex)
}

func endsWith(val string, regex string) bool{
	return strings.HasSuffix(val, regex)
}

func contains(val string, regex string) bool{
	return strings.Contains(val, regex)
}

func equals(val string, regex string) bool{
	return val == regex
}

func call(funcName string, params ... interface{}) (res []reflect.Value, err error) {
	defer func() {
        if err := recover(); err != nil {
			fmt.Println("panic occurred at formatter call function:", err)
        }
    }()
	f := reflect.ValueOf(funcMapper[funcName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is out of index.")
		return nil, err
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		
		in[k] = reflect.ValueOf(param)
	}
	res = f.Call(in)
	return res, nil
}