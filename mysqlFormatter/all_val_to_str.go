package mysqlFormatter

import(
	"fmt"
)

func (after AfterType) strConversion(tableName string, operation string, field string, finalPayload map[string]interface{}){
	if !canDoAllToStrD(tableName,operation, allValToStrOperations) {
		return
	}
	finalPayload[field] = fmt.Sprintf("%v", finalPayload[field])
}

func (before BeforeType) strConversion(tableName string, operation string, field string, finalPayload map[string]interface{}){
	if !canDoAllToStrD(tableName,operation, allValToStrOperations) {
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