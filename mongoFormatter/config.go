package mongoFormatter
import(
	"encoding/json"
	"errors"
	"regexp"
	"fmt"
	"strings"
)

var configMap = make(map[string]map[string]map[string]interface{})
var allRequiredAfterFields = make(map[string]map[string][]string)
var allRequiredBeforeFields = make(map[string]map[string][]string)
var allUnderscoreWithHostidFields = make(map[string]map[string]map[string]string)
var allUnderscoreWithHostidArrayFields = make(map[string]map[string]map[string]string)
var allNullToCollectionOperations = make(map[string]map[string]string)
var allValToStrOperations = make(map[string]map[string]bool)
var afterAllCollectionOperations = make(map[string]map[string]bool)
var collectionRegex = make(map[string]map[string]string)
var hostIdWithCollectioOperation = make(map[string]map[string]string)

var allowedKeys = map[string]bool{
	"insert": true,
	"update": true,
	"delete": true,
	"regex": true,
	"required_after_all": true,
	"required_after_fields": true,
	"required_before_fields": true,
	"underscore_with_hostid_fields": true,
	"underscore_with_hostid_array_fields": true,
	"all_value_to_string": true,
	"null_to": true,
	"host_id":true,
}

var stringBooleanKeys = map[string]string {
	"required_after_all": "true",
	"all_value_to_string": "true",
}

var commaSeperatedkeys = map[string]bool {
	"required_after_fields": true,
	"required_before_fields": true,
	"underscore_with_hostid_fields": true,
	"underscore_with_hostid_array_fields": true,
}

var splitToSliceFieldsFrmString = map[string]map[string]map[string][]string{
	"required_after_fields": allRequiredAfterFields,
	"required_before_fields": allRequiredBeforeFields,
}

var splitToMapFieldsFrmString = map[string]map[string]map[string]map[string]string{
	"underscore_with_hostid_fields": allUnderscoreWithHostidFields,
	"underscore_with_hostid_array_fields": allUnderscoreWithHostidArrayFields,
}


var strValuedKeys = map[string]bool {
	"required_after_all": true,
	"required_after_fields": true,
	"required_before_fields": true,
	"underscore_with_hostid_fields": true,
	"underscore_with_hostid_array_fields": true,
	"all_value_to_string": true,
	"null_to": true,
	"host_id": true,
}

var singleValuedKeys = map[string]string{
	"required_before_fields": "_id",
}

func Register(config string) error {
	if err := unMarshallConfig(config); err != nil {
		return errors.New(fmt.Sprintf("Invalid JSON Format for Config: %v", err))
	}

	if err := validateConfig(); err != nil {
		return err
	}

	return nil
}

func unMarshallConfig(config string) error{
	if err := json.Unmarshal([]byte(config), &configMap); err != nil {
		return err
	}
	return nil
}

func validateConfig() error{
	emptyConfig := true
	for collectionName,configVal := range configMap {
		emptyConfig = false
		if err := isValidOperationExist(collectionName, configMap[collectionName]); err != nil {
			return err
		}
		if err := isRegexExist(collectionName, configMap[collectionName]); err != nil {
			return err
		}
		for operation, v := range configVal {
			if isNotAnOperation(operation) {
				continue
			}
			if err := requiredAfterAllAndRequiredAfterFieldsCantBeTogether(v); err != nil {
				return err
			}
			if err := isValidKey(operation, operation); err != nil {
				return err
			}
			for k, val := range v {
				if err := isValidKey(operation, k); err != nil {
					return err
				}
				if err := onlyStringValuesAllowed(collectionName, operation, k, val); err != nil {
					return err
				}
			}
		}
	}
	if emptyConfig {
		return errors.New("You cant have empty configuration")
	}
	// fmt.Println("allRequiredAfterFields:", allRequiredAfterFields)
	// fmt.Println("afterAllCollectionOperations:", afterAllCollectionOperations)
	// fmt.Println("allRequiredBeforeFields:", allRequiredBeforeFields)
	// fmt.Println("allUnderscoreWithHostidFields", allUnderscoreWithHostidFields)
	// fmt.Println("allUnderscoreWithHostidArrayFields", allUnderscoreWithHostidArrayFields)
	// fmt.Println("allNullToCollectionOperations", allNullToCollectionOperations)
	// fmt.Println("allValToStrOperations", allValToStrOperations)
	return nil
}

func isNotAnOperation(operation string) bool{
	if operation != "insert" && operation != "update" && operation != "delete" {
		return true
	}
	return false
}

func isRegexExist(collectionName string, dataMap map[string]map[string]interface{}) error{
	if _,ok := dataMap["regex"]; ok {
		if len(dataMap["regex"]) != 1 {
			return errors.New("object regex only can have one key")
		}
		for k,v := range dataMap["regex"]{
			val, ok := canBeAssertToString(v)
			if !ok {
				return errors.New("object regex key value must be string")
			}
			if val != "startsWith" &&  val != "endsWith" && val != "contains" && val != "equals" {
				return errors.New("object regex key's value must be startsWith|endsWith|contains|equals")
			}
			if _, ok := collectionRegex[collectionName]; !ok {
				collectionRegex[collectionName] = make(map[string]string)
			}
			collectionRegex[collectionName][k] = val
		}
	}
	return nil
}

func isValidOperationExist(collectionName string, operations map[string]map[string]interface{}) error{
	_, insertOk := operations["insert"]
	_, updateOk := operations["update"]
	_, deleteOk := operations["delete"]
	if !insertOk && !updateOk && !deleteOk {
		return errors.New("collection "+collectionName+ " need atlest insert|update|delete operation to be mentioned")
	}
	return nil
} 

func isValidKey(operation string, key string) error {
	if _, ok := allowedKeys[key]; !ok {
		return errors.New("key "+key+ " is not allowed in configurations")
	}
	if key == "required_before_fields" && operation == "insert" {
		return errors.New("key "+key+ " is not allowed for operation "+operation)
	}
	if key == "required_after_all" && operation != "insert" {
		return errors.New("key "+key+ " is not allowed for operation "+operation)
	}
	if key == "required_after_fields" && operation != "insert" {
		return errors.New("key "+key+ " is not allowed for operation "+operation)
	}
	return nil
}

func requiredAfterAllAndRequiredAfterFieldsCantBeTogether(v map[string]interface{}) error {
	_, rqAllok := v["required_after_all"]
	_, rqFieldok := v["required_after_fields"]

	if rqAllok && rqFieldok {
		return errors.New("can't use required_after_fields while using required_after_all please remove required_after_fields")
	}
	if rqFieldok && rqAllok {
		return errors.New("can't use required_after_all while using required_after_fields please remove required_after_all")
	}
	return nil
}

func onlyStringValuesAllowed(collection string, operation string, key string, v interface{}) error {
	if _, ok := strValuedKeys[key]; !ok {
		return nil
	}
	val, ok := v.(string)
	if !ok {
		return errors.New(key+" value should be string")
	}
	if err := emptyValuesNotAllowed(collection, operation, key, val); err != nil {
		return err
	}
	if err := isAStrBoolKey(collection, operation, key, val); err != nil {
		return err
	}
	if err := isSingleValuedKey(key, val); err != nil {
		return err
	}
	ifNullToAvailable(collection, operation, key, val)
	ifHostIDAvailable(collection, operation, key, val)
	return nil
}

func isSingleValuedKey(key string, val string) error{
	if value, ok := singleValuedKeys[key]; ok {
		if val != value {
			return errors.New(key+" value should be "+value)
		}
	}
	return nil
}

func emptyValuesNotAllowed(collection string, operation string, key string, val string) error {
	if strings.TrimSpace(val) == "" {
		return errors.New(key+" value cant be empty")
	}
	if err := isACommaSeperatedValue(collection, operation, key, val); err != nil {
		return err
	}
	return nil
}

func isAStrBoolKey(collection string, operation string, key string, val string) error {
	if v, ok := stringBooleanKeys[key]; ok {
		if v != val {
			return errors.New(key+" value should be "+v)
		}
		setBoolCollectionOperation(collection, operation, key)
	}
	return nil
}

func isACommaSeperatedValue(collection string, operation string, key string, val string) error{
	r := regexp.MustCompile(`^(\w+)(,\s*\w+)*$`)
	if _, ok := commaSeperatedkeys[key]; ok {
		if !r.MatchString(val) {
			return errors.New(key+" value should be comma seperated")
		}
		if varCS, ok := splitToSliceFieldsFrmString[key]; ok {
			if _, ok := varCS[collection]; !ok {
				varCS[collection] = make(map[string][]string)
			}
			varCS[collection][operation] = strings.Split(val, ",")
			
		}
		if varCS, ok := splitToMapFieldsFrmString[key]; ok {
			var hostFieldsMap = make(map[string]string)
			for _,v := range strings.Split(val, ",") {
				hostFieldsMap[v] = v
			}
			if _, ok := varCS[collection]; !ok {
				varCS[collection] = make(map[string]map[string]string)
			}
			varCS[collection][operation] = hostFieldsMap
		}
	} 
	return nil
}

func ifHostIDAvailable(collection string, operation string, key string, v string) {
	if key != "host_id" {
		return
	}
	if _,ok := hostIdWithCollectioOperation[collection]; !ok {
		hostIdWithCollectioOperation[collection] = make(map[string]string)
	}
	hostIdWithCollectioOperation[collection][operation] = v
}

func ifNullToAvailable(collection string, operation string, key string, v string) {
	if key != "null_to" {
		return
	}
	if _,ok := allNullToCollectionOperations[collection]; !ok {
		allNullToCollectionOperations[collection] = make(map[string]string)
	}
	allNullToCollectionOperations[collection][operation] = v
}

func setBoolCollectionOperation(collection string, operation string, key string) {
	if key == "all_value_to_string" {
		if _,ok := allValToStrOperations[collection]; !ok {
			allValToStrOperations[collection] = make(map[string]bool)
		}
		allValToStrOperations[collection][operation] = true
	}
	if key == "required_after_all" {
		if _,ok := afterAllCollectionOperations[collection]; !ok {
			afterAllCollectionOperations[collection] = make(map[string]bool)
		}
		afterAllCollectionOperations[collection][operation] = true
	}
}






