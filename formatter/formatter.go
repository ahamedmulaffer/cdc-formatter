package formatter

import(
	"encoding/json"
	"github.com/ahamedmulaffer/cdc-formatter/mysqlFormatter"
	"reflect"
	// "fmt"
	"errors"
	"github.com/ahamedmulaffer/cdc-formatter/mongoFormatter"
)

var funcMapper = map[string]interface{}{
	"mysql": mysqlFormatter.Process,
	"mongodb": mongoFormatter.Process,
}

func Process(message string) (interface{}, error) {
	currentMessageMap := make(map[string]interface{})
	// currentMessagePayloadMap := make(map[string]interface{})
	currentMessageSourceMap := make(map[string]interface{})
	if err := messageUnMarshaller(message, currentMessageMap); err != nil {
		return nil, err
	}
	// payloadStr, _ := json.Marshal(currentMessageMap["payload"])
	// if err := messageUnMarshaller(string(payloadStr), currentMessagePayloadMap); err != nil {
	// 	return nil, err
	// }
	sourceStr, _ := json.Marshal(currentMessageMap["source"])
	if err := messageUnMarshaller(string(sourceStr), currentMessageSourceMap); err != nil {
		return nil, err
	}
	connectorStr := currentMessageSourceMap["connector"].(string)
	res, err := call(connectorStr, currentMessageMap, currentMessageSourceMap)
	if err != nil {
		return nil, err
	}
	return res[0].Interface().(map[string]map[string]interface{}), err
}

func call(funcName string, params ... interface{}) (res []reflect.Value, err error) {
	// defer func() {
    //     if err := recover(); err != nil {
	// 		fmt.Println("panic occurred at formatter call function:", err)
    //     }
    // }()
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

func messageUnMarshaller(message string, messageMap map[string]interface{}) error{
	if err := json.Unmarshal([]byte(message), &messageMap); err != nil {
		// fmt.Println("Error in Unmarshaling at currentMessage:", err)
		return err
	}
	return nil
}



