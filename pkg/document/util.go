package document

import (
	"fmt"
)

type jsonableMap map[string]interface{}

func convertMapKeyToString(key interface{}) string {
	switch key.(type) {
	case string:
		return key.(string)
	case int:
		return fmt.Sprintf("int(%d)", key)
	case float64:
		return fmt.Sprintf("float(%f)", key)
	case bool:
		return fmt.Sprintf("bool(%t)", key)
	}

	return fmt.Sprintf("?(%+v)", key)
}

// The json library can only marshal maps with string keys, and so all of our lists and maps that go into documentation
// must be converted to have only string keys before marshalling
func convertHelmValuesToJsonable(values interface{}) interface{} {
	switch values.(type) {
	case map[interface{}]interface{}:
		convertedMap := make(jsonableMap)

		for key, value := range values.(map[interface{}]interface{}) {
			convertedMap[convertMapKeyToString(key)] = convertHelmValuesToJsonable(value)
		}

		return convertedMap

	case []interface{}:
		convertedList := make([]interface{}, 0)

		for _, value := range values.([]interface{}) {
			convertedList = append(convertedList, convertHelmValuesToJsonable(value))
		}

		return convertedList

	default:
		return values
	}
}
