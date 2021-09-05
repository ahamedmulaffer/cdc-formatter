package mongoFormatter

func extract_id(val interface{}) interface{}{
	return val.(map[string]interface{})["$oid"]
}