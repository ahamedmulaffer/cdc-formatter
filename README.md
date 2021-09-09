# Welcome to CDC-Formatter

cdc-formatter is a library to extract needed fields and do some conversion efficiently from debezium cdc json response while keeping the debezium cdc engine as common

currently it is implemented only for mysql and mongo only

## How to use it
go get github.com/ahamedmulaffer/cdc-formatter

## Example
```go
package main

import(
    "github.com/ahamedmulaffer/cdc-formatter/mysqlFormatter"
    "github.com/ahamedmulaffer/cdc-formatter/mongoFormatter"
    "fmt"
)
func main(){
    // here table is the name of a mysql table
    // insert,update.delete are the only permitted operations or keys inside a table object
    // required_after_fields,required_before_fields should be a comma seperated fields it will only fetch the mentioned fields
    // NOTE :- you cant use required_after_all and required_after_fields together in one operation either you have to use one of them
    //and required_before_all and required_before_fields cant be together in one operation 
     
    // time_formatter is currenttly supported only for mysql and supported data types are time,date,datetime
    // here field3 is a time field in mysql and we need the time to be in hh:mm:ss format
    // thats why its used as TIME|15:04:05 here TIME is the actual data type of the field 15:04:05 is the golang layout format for hh:mm:ss you can use any valid golang format to convert it
    // all_value_to_string will set all values to string
    mysqlConfig := `{
        "table":{
            "insert":{
                "required_after_all": "true",
                "required_after_fields": "field1,field2",
                "time_formatter": {
                    "field3": "TIME|15:04:05",
                    "field4": "DATE|2006-01-02",
                    "field5": "DATETIME|2006-01-02 15:04:05"
                },
                "all_value_to_string": "true"
            },
            "update":{
                "required_after_all": "true",
                "required_before_all": "true",
                "required_before_fields": "field1, field2"
                "time_formatter": {
                    "field3": "TIME|15:04:05",
                    "field4": "DATE|2006-01-02",
                    "field5": "DATETIME|2006-01-02 15:04:05"
                },
                "all_value_to_string": "true"
            },
            "delete":{
                "required_before_fields": "field1, field2"
            }
        }
    }`

    err := mysqlFormatter.Register(mysqlConfig)
    if err != nil {
        fmt.Println(err)
        return
    }

    mongoConfig :=`{
        "collectionname": {
            "insert":{
                "required_after_all": "true",
                "required_after_fields": "field1,field2",
                "all_value_to_string": "true"
            },
            "update":{
                "required_before_fields": "_id",
                "all_value_to_string": "true"
            },
            "delete":{
                "required_before_fields": "_id",
                "all_value_to_string": "true"
            }
        }
    }`

    /*
    mongoConfig :=`{
        "collectionname": {
            "regex": {
                    "value": "equals" // here any collection name equal to value will be mapped to
                                        collectionname supprted regex are equals,startsWith,endsWith,contains
            },
            "insert":{
                    "required_after_all": "true",
                    "required_after_fields": "field1,field2",
                    "all_value_to_string": "true"
                },
                "update":{
                    "required_before_fields": "_id",
                    "all_value_to_string": "true"
                },
                "delete":{
                    "required_before_fields": "_id",
                    "all_value_to_string": "true"
                }
        }
    }`
    */
    err := mongoFormatter.Register(mongoConfig)
    if err != nil {
        fmt.Println(err)
        return
    }
    //response, err := formatter.Process(response from debezium cdc)
}

```