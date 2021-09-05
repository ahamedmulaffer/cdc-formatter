package mysqlFormatter

import(
	"strings"
	"fmt"
)

func canBeAssertToString(val interface{}) (string, bool) {
	value, ok := val.(string)
	if ok {
		return strings.TrimSpace(value) ,ok
	}
	return "", ok
}

func assertToString(val interface{}) string {
   var strVal string
   var ok bool
   strVal, ok = canBeAssertToString(val)
   if ok {
	   return strVal
   }
   var float64Val float64
   float64Val, ok = val.(float64)
   if ok {
		strVal = fmt.Sprintf("%.0f", float64Val)
		return strVal
   }
   return strVal
}