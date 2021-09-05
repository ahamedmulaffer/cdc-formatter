package mongoFormatter

import(
	"fmt"
)

func (after AfterType) strConversion(collectionName string, operation string, field string, finalPayload map[string]interface{}){
	if !canDoAllToStrD(collectionName,operation, allValToStrOperations) {
		return
	}
	finalPayload[field] = fmt.Sprintf("%v", finalPayload[field])
}

func (before BeforeType) strConversion(collectionName string, operation string, field string, finalPayload map[string]interface{}){
	if !canDoAllToStrD(collectionName,operation, allValToStrOperations) {
		return
	}
	finalPayload[field] = fmt.Sprintf("%v", finalPayload[field])
}

func canDoAllToStrD(collectionName string, operation string, dataMap map[string]map[string]bool) bool{
	if len(dataMap) == 0 {
		return false
	}
	if _, ok := dataMap[collectionName]; !ok {
		return false
	}
	if len(dataMap[collectionName]) == 0 {
		return false
	}
	if _, ok := dataMap[collectionName][operation]; !ok {
		return false
	}
	return true
}