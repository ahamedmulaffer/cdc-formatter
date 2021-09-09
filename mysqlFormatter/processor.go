package mysqlFormatter

import(
	// "strings"
	// "fmt"
)

var allowedOperations = map[string]string{
	"c": "insert",
	"u": "update",
	"d": "delete",
}

func Process(payload map[string]interface{}, source map[string]interface{}) map[string]map[string]interface{}{
	finalPayload := make(map[string]map[string]interface{})
	tableName := "N/A"
	operation := "N/A"
	finalPayload["before"] = make(map[string]interface{})
	finalPayload["after"] = make(map[string]interface{})
	finalPayload["source"] = map[string]interface{} {
		"database": "mysql",
		"table": tableName,
		"operation": operation,
		"allowed": false,
	}
	if _, ok := payload["op"]; !ok {
		return finalPayload
	}
	var afterType AfterType
	var beforeType BeforeType
	tableName = getTableName(source)
	operation = allowedOperations[payload["op"].(string)]
	finalPayload["source"]["table"] = tableName
	finalPayload["source"]["operation"] = operation
	//check table and its operation registered
	if _, ok := configMap[tableName]; !ok {
		return finalPayload
	}
	if _, ok := configMap[tableName][operation]; !ok {
		return finalPayload
	}
	finalPayload["source"]["allowed"] = true
	requiredAfterFields , requiredAfterFieldsAvailable := afterType.isRequiredFieldsAvailable(tableName, operation)
	requiredB4Fields , requiredB4FieldsAvailable := beforeType.isRequiredFieldsAvailable(tableName, operation)
	if payload["after"] != nil {
		if requiredAfterFieldsAvailable{
			afterType.loopRequiredFields(tableName,operation,requiredAfterFields,payload["after"].(map[string]interface{}),finalPayload["after"])
		}
		if _, ok := afterAllTableOperations[tableName][operation]; ok {
			afterType.loopPayloadFields(tableName,operation,payload["after"].(map[string]interface{}),finalPayload["after"])
		}
		
	}
	//Before
	if payload["before"] != nil {
		if requiredB4FieldsAvailable {
			beforeType.loopRequiredFields(tableName,operation,requiredB4Fields,payload["before"].(map[string]interface{}),finalPayload["before"])
		}
		if _, ok := beforeAllTableOperations[tableName][operation]; ok {
			beforeType.loopPayloadFields(tableName,operation,payload["before"].(map[string]interface{}),finalPayload["before"])
		}
	}
	// fmt.Println("finalPayload", finalPayload)
	return finalPayload
}

func getTableName(source map[string]interface{}) string{
	return source["table"].(string)
}
