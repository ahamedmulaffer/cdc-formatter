package mysqlFormatter

import(
	// "fmt"
)

type BeforeType string



func (before BeforeType) loopRequiredFields(tableName string, operation string, requiredB4Fields []string, payloadBefore map[string]interface{}, finalPayload map[string]interface{}){
	idToHPIDFuncCalled := false
	for _, reqField := range requiredB4Fields {
		if idToHPIDFuncCalled && reqField == "HPID" {
			continue
		}
		//check if map to id enbled
		if sliceField(requiredB4Fields).isMapIDToHPIDEnabled(tableName, operation, idToHPIDFuncCalled) && reqField == "ID"{
			idToHPIDFuncCalled = true
			if _, ok := payloadBefore["ID"]; ok {
				finalPayload["HPID"] = payloadBefore["ID"]
				delete(payloadBefore, "ID")
			}
			continue
			
		}
		//check isNullToAvailable if value nil only
		if payloadBefore[reqField] == nil {
			val, ok := isNullToAvailable(tableName, operation)
			if ok {
				finalPayload[reqField] = val
			}
			continue
		}
		finalPayload[reqField] = payloadBefore[reqField]
		if _, ok:= payloadBefore["HostID"]; ok {
			hostId := assertToString(payloadBefore["HostID"])
			before.underscoreWithHostidFields(tableName,operation,reqField,finalPayload,hostId)
			before.underscoreWithHostidAssetFields(tableName,operation,reqField,finalPayload,hostId)
		}
		before.timeFieldConversion(tableName,operation,reqField,finalPayload)
		before.strConversion(tableName,operation,reqField,finalPayload)

	}
}

func (before BeforeType) isRequiredFieldsAvailable(tableName string, operation string) ([]string, bool){
	if len(allRequiredBeforeFields) == 0 {
		return nil, false
	}
	if _, ok := allRequiredBeforeFields[tableName]; !ok {
		return nil, false
	}
	if len(allRequiredBeforeFields[tableName]) == 0 {
		return nil, false
	}
	if _, ok := allRequiredBeforeFields[tableName][operation]; !ok {
		return nil, false
	}
	if len(allRequiredBeforeFields[tableName][operation]) == 0 {
		return nil, false
	}
	return allRequiredBeforeFields[tableName][operation], true
}

func (before BeforeType) loopPayloadFields(tableName string, operation string, payloadBefore map[string]interface{}, finalPayload map[string]interface{}) {
	idToHPIDFuncCalled := false
	for field, _ := range payloadBefore {
		if idToHPIDFuncCalled && field == "HPID" {
			continue
		}
		//check if map to id enbled
		if mapField(payloadBefore).isMapIDToHPIDEnabled(tableName, operation, idToHPIDFuncCalled) && field == "ID"{
			idToHPIDFuncCalled = true
			if _, ok := payloadBefore["ID"]; ok {
				finalPayload["HPID"] = payloadBefore["ID"]
				delete(payloadBefore, "ID")
			}
			continue
		}
		//check isNullToAvailable if value nil only
		if payloadBefore[field] == nil {
			val, ok := isNullToAvailable(tableName, operation)
			if ok {
				finalPayload[field] = val
			}
			continue
		}
		finalPayload[field] = payloadBefore[field]
		if _, ok:= payloadBefore["HostID"]; ok {
			hostId := assertToString(payloadBefore["HostID"])
			before.underscoreWithHostidFields(tableName,operation,field,finalPayload,hostId)
			before.underscoreWithHostidAssetFields(tableName,operation,field,finalPayload,hostId)
		}
		if val, ok := isCreateOrUpdateByAvailable(tableName, operation, field); ok {
			finalPayload[field] = val
		}
		before.timeFieldConversion(tableName,operation,field,finalPayload)
		before.strConversion(tableName,operation,field,finalPayload)
	}
}