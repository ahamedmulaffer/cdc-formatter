package mysqlFormatter
import(
	"encoding/json"
	"errors"
	"regexp"
	"fmt"
	"strings"
)

var configMap = make(map[string]map[string]map[string]interface{})
var allTimeFormatterMap = make(map[string]map[string]map[string]string)
var allRequiredAfterFields = make(map[string]map[string][]string)
var allRequiredBeforeFields = make(map[string]map[string][]string)
var allUnderscoreWithHostidFields = make(map[string]map[string]map[string]string)
var allUnderscoreWithHostidAssetFields = make(map[string]map[string]map[string]string)
var allNullToTableOperations = make(map[string]map[string]string)
var allIDToHPIDTableOperations = make(map[string]map[string]bool)
var allValToStrOperations = make(map[string]map[string]bool)
var afterAllTableOperations = make(map[string]map[string]bool)
var beforeAllTableOperations = make(map[string]map[string]bool)

var allowedKeys = map[string]bool{
	"insert": true,
	"update": true,
	"delete": true,
	"required_after_all": true,
	"required_before_all": true,
	"required_after_fields": true,
	"required_before_fields": true,
	"underscore_with_hostid_fields": true,
	"underscore_with_hostid_asset_fields": true,
	"time_formatter": true,
	"all_value_to_string": true,
	"null_to": true,
	"map_id_to_hpid": true,
}

var stringBooleanKeys = map[string]string {
	"required_after_all": "true",
	"required_before_all": "true",
	"all_value_to_string": "true",
	"map_id_to_hpid": "true",
}

var commaSeperatedkeys = map[string]bool {
	"required_after_fields": true,
	"required_before_fields": true,
	"underscore_with_hostid_fields": true,
	"underscore_with_hostid_asset_fields": true,
}

var splitToSliceFieldsFrmString = map[string]map[string]map[string][]string{
	"required_after_fields": allRequiredAfterFields,
	"required_before_fields": allRequiredBeforeFields,
}

var splitToMapFieldsFrmString = map[string]map[string]map[string]map[string]string{
	"underscore_with_hostid_fields": allUnderscoreWithHostidFields,
	"underscore_with_hostid_asset_fields": allUnderscoreWithHostidAssetFields,
}


var strValuedKeys = map[string]bool {
	"required_after_all": true,
	"required_before_all": true,
	"required_after_fields": true,
	"required_before_fields": true,
	"underscore_with_hostid_fields": true,
	"underscore_with_hostid_asset_fields": true,
	"all_value_to_string": true,
	"null_to": true,
	"map_id_to_hpid": true,
}
var objectKeys = map[string]bool {
	"time_formatter": true,
}

var timeValueFormat = map[string]string{
	"TIME": "TIME",
	"DATE": "DATE",
	"DATETIME": "DATETIME",
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
	for tableName,configVal := range configMap {
		emptyConfig = false
		if err := isValidOperationExist(tableName, configMap[tableName]); err != nil {
			return err
		}
		for operation, v := range configVal {
			if err := requiredAfterAllAndRequiredAfterFieldsCantBeTogether(v); err != nil {
				return err
			}
			if err := requiredBeforeAllAndRequiredBeforeFieldsCantBeTogether(v); err != nil {
				return err
			}
			if err := isValidKey(operation, operation); err != nil {
				return err
			}
			for k, val := range v {
				if err := isValidKey(operation, k); err != nil {
					return err
				}
				if err := onlyStringValuesAllowed(tableName, operation, k, val); err != nil {
					return err
				}
				if err := isTimeFormatterValid(tableName, operation, k, val); err != nil {
					return err
				}
			}
		}
	}
	if emptyConfig {
		return errors.New("You cant have empty configuration")
	}
	// fmt.Println("allRequiredAfterFields:", allRequiredAfterFields)
	// fmt.Println("allRequiredBeforeFields:", allRequiredBeforeFields)
	// fmt.Println("allUnderscoreWithHostidFields", allUnderscoreWithHostidFields)
	// fmt.Println("allUnderscoreWithHostidAssetFields", allUnderscoreWithHostidAssetFields)
	// fmt.Println("allTimeFormatterMap", allTimeFormatterMap)
	// fmt.Println("allNullToTableOperations", allNullToTableOperations)
	// fmt.Println("allIDToHPIDTableOperations", allIDToHPIDTableOperations)
	// fmt.Println("allValToStrOperations", allValToStrOperations)
	return nil
}

func isValidOperationExist(tableName string, operations map[string]map[string]interface{}) error{
	_, insertOk := operations["insert"]
	_, updateOk := operations["update"]
	_, deleteOk := operations["delete"]
	if !insertOk && !updateOk && !deleteOk {
		return errors.New("table "+tableName+ " need atlest insert|update|delete operation to be mentioned")
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
	if key == "required_before_all" && operation == "insert" {
		return errors.New("key "+key+ " is not allowed for operation "+operation)
	}
	if key == "required_after_all" && operation == "delete" {
		return errors.New("key "+key+ " is not allowed for operation "+operation)
	}
	if key == "required_after_fields" && operation == "delete" {
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

func requiredBeforeAllAndRequiredBeforeFieldsCantBeTogether(v map[string]interface{}) error {
	_, rqAllok := v["required_before_all"]
	_, rqFieldok := v["required_before_fields"]

	if rqAllok && rqFieldok {
		return errors.New("can't use required_before_fields while using required_before_all please remove required_before_fields")
	}
	if rqFieldok && rqAllok {
		return errors.New("can't use required_before_all while using required_before_fields please remove required_before_all")
	}
	return nil
}

func onlyStringValuesAllowed(table string, operation string, key string, v interface{}) error {
	if _, ok := strValuedKeys[key]; !ok {
		return nil
	}
	val, ok := v.(string)
	if !ok {
		return errors.New(key+" value should be string")
	}
	if err := emptyValuesNotAllowed(table, operation, key, val); err != nil {
		return err
	}
	if err := isAStrBoolKey(table, operation, key, val); err != nil {
		return err
	}
	ifNullToAvailable(table, operation, key, val)
	return nil
}

func emptyValuesNotAllowed(table string, operation string, key string, val string) error {
	if strings.TrimSpace(val) == "" {
		return errors.New(key+" value cant be empty")
	}
	if err := isACommaSeperatedValue(table, operation, key, val); err != nil {
		return err
	}
	return nil
}

func isAStrBoolKey(table string, operation string, key string, val string) error {
	if v, ok := stringBooleanKeys[key]; ok {
		if v != val {
			return errors.New(key+" value should be "+v)
		}
		setBoolTableOperation(table, operation, key)
	}
	return nil
}

func isACommaSeperatedValue(table string, operation string, key string, val string) error{
	r := regexp.MustCompile(`^(\w+)(,\s*\w+)*$`)
	if _, ok := commaSeperatedkeys[key]; ok {
		if !r.MatchString(val) {
			return errors.New(key+" value should be comma seperated")
		}
		if varCS, ok := splitToSliceFieldsFrmString[key]; ok {
			if _, ok := varCS[table]; !ok {
				varCS[table] = make(map[string][]string)
			}
			varCS[table][operation] = strings.Split(val, ",")
			
		}
		if varCS, ok := splitToMapFieldsFrmString[key]; ok {
			var hostFieldsMap = make(map[string]string)
			for _,v := range strings.Split(val, ",") {
				hostFieldsMap[v] = v
			}
			if _, ok := varCS[table]; !ok {
				varCS[table] = make(map[string]map[string]string)
			}
			varCS[table][operation] = hostFieldsMap
		}
	} 
	return nil
}

func isTimeFormatterValid(table string, operation string, key string, v interface{}) error{
	if _, ok := objectKeys[key]; ok {
		timeFormatStr, _ := json.Marshal(v)
		timeFormatterMap := make(map[string]string)
		if err := json.Unmarshal([]byte(timeFormatStr), &timeFormatterMap); err != nil {
			return errors.New(key+" value should be object of string key and string value")
		}
		for _, val := range timeFormatterMap{
			val = strings.TrimSpace(val)
			slice := strings.Split(val, "|")
			if len(slice) != 2{
				return errors.New(key+" value should be A|B format")
			}
			if _,ok := timeValueFormat[slice[0]]; !ok {
				return errors.New(key+"'s TimeKey "+slice[0]+" is not valid")
			}
		}
		if _, ok := allTimeFormatterMap[table]; !ok {
			allTimeFormatterMap[table] = make(map[string]map[string]string)
		}
		allTimeFormatterMap[table][operation] = timeFormatterMap
	}
	return nil
}

func ifNullToAvailable(table string, operation string, key string, v string) {
	if key != "null_to" {
		return
	}
	if _,ok := allNullToTableOperations[table]; !ok {
		allNullToTableOperations[table] = make(map[string]string)
	}
	allNullToTableOperations[table][operation] = v
}

func setBoolTableOperation(table string, operation string, key string) {
	if key == "all_value_to_string" {
		if _,ok := allValToStrOperations[table]; !ok {
			allValToStrOperations[table] = make(map[string]bool)
		}
		allValToStrOperations[table][operation] = true
	}
	if key == "map_id_to_hpid" {
		if _,ok := allIDToHPIDTableOperations[table]; !ok {
			allIDToHPIDTableOperations[table] = make(map[string]bool)
		}
		allIDToHPIDTableOperations[table][operation] = true
	}
	if key == "required_after_all" {
		if _,ok := afterAllTableOperations[table]; !ok {
			afterAllTableOperations[table] = make(map[string]bool)
		}
		afterAllTableOperations[table][operation] = true
	}
	if key == "required_before_all" {
		if _,ok := beforeAllTableOperations[table]; !ok {
			beforeAllTableOperations[table] = make(map[string]bool)
		}
		beforeAllTableOperations[table][operation] = true
	}

}






