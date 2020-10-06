package document

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	AlphaNumSortOrder = "alphanum"
	FileSortOrder     = "file"
)

// The json library can only marshal maps with string keys, and so all of our lists and maps that go into documentation
// must be converted to have only string keys before marshalling
func convertHelmValuesToJsonable(values *yaml.Node) interface{} {
	switch values.Kind {
	case yaml.MappingNode:
		convertedMap := make(map[string]interface{})

		for i := 0; i < len(values.Content); i += 2 {
			k := values.Content[i]
			v := values.Content[i+1]
			convertedMap[k.Value] = convertHelmValuesToJsonable(v)
		}

		return convertedMap
	case yaml.SequenceNode:
		convertedList := make([]interface{}, 0)

		for _, v := range values.Content {
			convertedList = append(convertedList, convertHelmValuesToJsonable(v))
		}

		return convertedList
	case yaml.AliasNode:
		return convertHelmValuesToJsonable(values.Alias)
	case yaml.ScalarNode:
		switch values.Tag {
		case nullTag:
			return nil
		case strTag:
			fallthrough
		case timestampTag:
			return values.Value
		case intTag:
			var decodedValue int
			err := values.Decode(&decodedValue)
			if err != nil {
				log.Errorf("Failed to decode value from yaml node value %s", values.Value)
				return 0
			}
			return decodedValue
		case floatTag:
			var decodedValue float64
			err := values.Decode(&decodedValue)
			if err != nil {
				log.Errorf("Failed to decode value from yaml node value %s", values.Value)
				return 0
			}
			return decodedValue

		case boolTag:
			var decodedValue bool
			err := values.Decode(&decodedValue)
			if err != nil {
				log.Errorf("Failed to decode value from yaml node value %s", values.Value)
				return 0
			}
			return decodedValue
		}
	}

	return nil
}
