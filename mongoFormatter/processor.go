package mongoFormatter

import(
	// "strings"
	// "fmt"
	"encoding/json"
)

var allowedOperations = map[string]string{
	"c": "insert",
	"u": "update",
	"d": "delete",
}

func Process(payload map[string]interface{}, source map[string]interface{}) map[string]map[string]interface{}{
	finalPayload := make(map[string]map[string]interface{})
	collectionName := "N/A"
	operation := "N/A"
	finalPayload["before"] = make(map[string]interface{})
	finalPayload["after"] = make(map[string]interface{})
	finalPayload["source"] = map[string]interface{} {
		"database": "mongo",
		"collection": collectionName,
		"operation": operation,
		"allowed": false,
	}
	if _, ok := payload["op"]; !ok {
		return finalPayload
	}
	var afterType AfterType
	var beforeType BeforeType
	collectionName = getCollectionName(source)
	operation = allowedOperations[payload["op"].(string)]
	finalPayload["source"]["collection"] = collectionName
	finalPayload["source"]["operation"] = operation
	
	//check collection and its operation registered
	if _, ok := configMap[collectionName]; !ok {
		return finalPayload
	}
	if _, ok := configMap[collectionName][operation]; !ok {
		return finalPayload
	}
	finalPayload["source"]["allowed"] = true
	requiredAfterFields , requiredAfterFieldsAvailable := afterType.isRequiredFieldsAvailable(collectionName, operation)
	requiredB4Fields , requiredB4FieldsAvailable := beforeType.isRequiredFieldsAvailable(collectionName, operation)
	if payload["after"] != nil {
		payloadAfter := make(map[string]interface{})
		if err := json.Unmarshal([]byte(payload["after"].(string)), &payloadAfter); err != nil {
			return finalPayload
		}
		if requiredAfterFieldsAvailable{
			afterType.loopRequiredFields(collectionName,operation,requiredAfterFields,payloadAfter,finalPayload["after"])
		}
		if _, ok := afterAllCollectionOperations[collectionName][operation]; ok {
			afterType.loopPayloadFields(collectionName,operation,payloadAfter,finalPayload["after"])
		}
		
	}
	//Before
	if payload["filter"] != nil {
		payloadBefore := make(map[string]interface{})
		if err := json.Unmarshal([]byte(payload["filter"].(string)), &payloadBefore); err != nil {
			return finalPayload
		}
		if requiredB4FieldsAvailable {
			beforeType.loopRequiredFields(collectionName,operation,requiredB4Fields,payloadBefore,finalPayload["before"])
		}
	}
	// fmt.Println("finalPayload", finalPayload)
	return finalPayload
}

func getCollectionName(source map[string]interface{}) string{
	collName := source["collection"].(string)
	_, ok := configMap[collName]
	if ok {
		return collName
	}
	for collectionName,v := range collectionRegex {
		for pattern, regex := range v {
			res, _ := call(regex, collName, pattern)
			if res[0].Interface().(bool) {
				return collectionName
			}
			return collName
		}
	}
	return collName
}

