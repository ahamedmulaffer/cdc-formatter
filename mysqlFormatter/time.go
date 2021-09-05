package mysqlFormatter

import(
	"time"
	"strings"
	"reflect"
)

func (after AfterType) timeFieldConversion(tableName string, operation string, field string, finalPayload map[string]interface{}) {
	format, ok := canDoTimeConversion(tableName,operation,field,allTimeFormatterMap)
	if !ok {
		return
	}
	slice := strings.Split(format, "|")
	call(after, slice[0], slice[1], field, finalPayload)
}

func (before BeforeType) timeFieldConversion(tableName string, operation string, field string, finalPayload map[string]interface{}) {
	format, ok := canDoTimeConversion(tableName,operation,field,allTimeFormatterMap)
	if !ok {
		return
	}
	slice := strings.Split(format, "|")
	call(before, slice[0], slice[1], field, finalPayload)
}

func canDoTimeConversion(tableName string, operation string, field string, dataMap map[string]map[string]map[string]string) (string, bool){
	if len(dataMap) == 0 {
		return "", false
	}
	if _, ok := dataMap[tableName]; !ok {
		return "", false
	}
	if len(dataMap[tableName]) == 0 {
		return "", false
	}
	if _, ok := dataMap[tableName][operation]; !ok {
		return "", false
	}
	if len(dataMap[tableName][operation]) == 0 {
		return "", false
	}
	format, ok := dataMap[tableName][operation][field]
	if !ok {
		return "", false
	}
	return format, true
}

func (after AfterType) TIME(layout string, field string, finalPayload map[string]interface{}){
	ms := finalPayload[field].(float64)/1000
	t := time.Unix(0, int64(ms)*int64(time.Millisecond))
	finalPayload[field] = t.Format(layout)
}

func (before BeforeType) TIME(layout string, field string, finalPayload map[string]interface{}){
	ms := finalPayload[field].(float64)/1000
	t := time.Unix(0, int64(ms)*int64(time.Millisecond))
	finalPayload[field] = t.Format(layout)
}

func (after AfterType) DATE(layout string, field string, finalPayload map[string]interface{}){
	t := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, finalPayload[field].(int))
	finalPayload[field] = t.Format(layout)
}

func (before BeforeType) DATE(layout string, field string, finalPayload map[string]interface{}){
	t := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, finalPayload["before_"+field].(int))
	finalPayload[field] = t.Format(layout)
}

func (after AfterType) DATETIME(layout string, field string, finalPayload map[string]interface{}){
	t := time.Unix(0, int64(finalPayload["before_"+field].(float64))*int64(time.Millisecond))
	finalPayload[field] = t.Format(layout)
}

func (before BeforeType) DATETIME(layout string, field string, finalPayload map[string]interface{}){
	t := time.Unix(0, int64(finalPayload[field].(float64))*int64(time.Millisecond))
	finalPayload[field] = t.Format(layout)
}

func call(any interface{}, funcName string, args ... interface{}){
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(funcName).Call(inputs)
}
