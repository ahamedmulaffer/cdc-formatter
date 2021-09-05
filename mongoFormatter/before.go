package mongoFormatter

type BeforeType string


func (before BeforeType) loopRequiredFields(collectionName string, operation string, requiredB4Fields []string, payloadBefore map[string]interface{}, finalPayload map[string]interface{}){
	for _, reqField := range requiredB4Fields {
		if reqField == "id" || reqField == "_id"{
			finalPayload["id"] = extract_id(payloadBefore["_id"])
			continue
		}
		//check isNullToAvailable if value nil only
		if payloadBefore[reqField] == nil {
			val, ok := isNullToAvailable(collectionName, operation)
			if ok {
				finalPayload[reqField] = val
			}
			continue
		}
		finalPayload[reqField] = payloadBefore[reqField]
		if hostId, ok:= hostIdAvailable(collectionName, operation); ok {
			before.underscoreWithHostidFields(collectionName,operation,reqField,finalPayload,hostId)
			before.underscoreWithHostidArrayFields(collectionName,operation,reqField,finalPayload,hostId)
		}
		before.strConversion(collectionName,operation,reqField,finalPayload)

	}
}

func (before BeforeType) isRequiredFieldsAvailable(collectionName string, operation string) ([]string, bool){
	if len(allRequiredBeforeFields) == 0 {
		return nil, false
	}
	if _, ok := allRequiredBeforeFields[collectionName]; !ok {
		return nil, false
	}
	if len(allRequiredBeforeFields[collectionName]) == 0 {
		return nil, false
	}
	if _, ok := allRequiredBeforeFields[collectionName][operation]; !ok {
		return nil, false
	}
	if len(allRequiredBeforeFields[collectionName][operation]) == 0 {
		return nil, false
	}
	return allRequiredBeforeFields[collectionName][operation], true
}

func (before BeforeType) loopPayloadFields(collectionName string, operation string, payloadBefore map[string]interface{}, finalPayload map[string]interface{}) {
	for field, _ := range payloadBefore {
		if field == "id" || field == "_id"{
			finalPayload["id"] = extract_id(payloadBefore["_id"])
			continue
		}
		//check isNullToAvailable if value nil only
		if payloadBefore[field] == nil {
			val, ok := isNullToAvailable(collectionName, operation)
			if ok {
				finalPayload[field] = val
			}
			continue
		}
		finalPayload[field] = payloadBefore[field]
		if hostId, ok:= hostIdAvailable(collectionName, operation); ok {
			before.underscoreWithHostidFields(collectionName,operation,field,finalPayload,hostId)
			before.underscoreWithHostidArrayFields(collectionName,operation,field,finalPayload,hostId)
		}
		before.strConversion(collectionName,operation,field,finalPayload)
	}
}