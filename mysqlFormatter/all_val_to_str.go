package mysqlFormatter

import(
	"fmt"
	"encoding/json"
	"reflect"
)

func (after AfterType) strConversion(tableName string, operation string, field string, finalPayload map[string]interface{}){
	if !canDoAllToStrD(tableName,operation, allValToStrOperations) {
		return
	}
	if isMapOrSlice(field, finalPayload) {
		val, _ := json.Marshal(finalPayload[field])
		finalPayload[field] = string(val)
		return
	}
	finalPayload[field] = fmt.Sprintf("%v", finalPayload[field])
}

func (before BeforeType) strConversion(tableName string, operation string, field string, finalPayload map[string]interface{}){
	if !canDoAllToStrD(tableName,operation, allValToStrOperations) {
		return
	}
	if isMapOrSlice(field, finalPayload) {
		val, _ := json.Marshal(finalPayload[field])
		finalPayload[field] = string(val)
		return
	}
	finalPayload[field] = fmt.Sprintf("%v", finalPayload[field])
}

func canDoAllToStrD(tableName string, operation string, dataMap map[string]map[string]bool) bool{
	if len(dataMap) == 0 {
		return false
	}
	if _, ok := dataMap[tableName]; !ok {
		return false
	}
	if len(dataMap[tableName]) == 0 {
		return false
	}
	if _, ok := dataMap[tableName][operation]; !ok {
		return false
	}
	return true
}

func isMapOrSlice(field string, finalPayload map[string]interface{}) bool{
	if reflect.ValueOf(finalPayload[field]).Kind() == reflect.Map {
		return true
	}

	if reflect.ValueOf(finalPayload[field]).Kind() == reflect.Slice {
		return true
	}

	return false
}