package matcher

import (
	"github.com/homedepot/github-webhook/configuration"
	"github.com/homedepot/github-webhook/events"
	"strings"
	"fmt"
	"reflect"
)

var debug bool

func addFlattenedData(flattened map[string][]string, prefix string, name string, value string) {
	if len(prefix) > 0 {
		name = prefix + "." + name
	}
	values := flattened[name]
	values = append(values,value)
	flattened[name] = values
}

func flattenData(flattened map[string][]string, prefix string, data map[string]interface{}) {
	for name, value := range data {
		switch value.(type) {
		case nil:
			addFlattenedData(flattened, prefix, name, "")
		case string, int, uint, uint8, uint16, uint32, bool, uint64, int8, int16, int32, int64, float32, float64, complex64, complex128:
			addFlattenedData(flattened, prefix, name, fmt.Sprint(value))
		case map[string]interface{}:
			m := value.(map[string]interface{})
			newPrefix := name
			if len(prefix) > 0 {
				newPrefix = prefix + "." + name
			}
			flattenData(flattened, newPrefix, m)
		case []string:
			addFlattenedData(flattened, prefix, name, strings.Join(value.([]string),","))
		case []interface{}:

			imap := value.([]interface{})
			if len(imap) == 0 {
				continue
			}
			val := value.([]interface{})[0]

			newPrefix := name
			if len(prefix) > 0 {
				newPrefix = prefix + "." + name
			}
			_, isMapSlice := val.([]map[string]interface{})
			if isMapSlice {
				for _, m := range imap {
					flattenData(flattened, newPrefix, m.(map[string]interface{}))
				}
			} else {
				_, isStringSlice := val.([]string)
				if isStringSlice {
					for _, m := range imap {
						addFlattenedData(flattened, prefix, name, strings.Join(m.([]string),","))
					}
				} else {
					m, isMap := val.(map[string]interface{})
					if isMap {
						flattenData(flattened, newPrefix, m)
					} else {
						s, isString := val.(string)
						if isString {
							addFlattenedData(flattened, prefix, name, s)
						} else {
							fmt.Println("1unsupported", reflect.TypeOf(val), reflect.TypeOf(value), name, value)
						}
					}
				}
			}

		default:
			fmt.Println("2unsupported",name,value)
		}

	}
}

func Matches(trigger configuration.Trigger, event events.Event) bool {

	if !strings.EqualFold(trigger.Event,event.Name) {
		return false
	}

	if len(trigger.Rules) == 0 {
		return true
	}

	flattened := make(map[string][]string)
	flattenData(flattened,"", event.Data)

	matches := true

	var matchFound bool

	if debug {
		fmt.Println("matching", event.Name)
	}

	for name, value := range trigger.Rules {
		matchFound = false
		if debug {
			fmt.Println("attempting to match rule",name)
		}
		for _, flatValue := range flattened[name] {
			if debug {
				fmt.Println("checking",name,value,"vs",flatValue)
			}
			if strings.EqualFold(value, flatValue) {
				matchFound = true
				break
			}
			if debug {
				if matchFound {
					fmt.Println("match found for",name)
				} else {
					fmt.Println("match NOT found for",name)
				}
			}
		}
		matches = matches && matchFound
	}

	if debug {
		if matchFound {
			fmt.Println("trigger matches",event.Name)
		} else {
			fmt.Println("NO match for trigger",event.Name)
		}
	}

	return matches
}

func SetDebug(d bool) {
	debug = d
}
