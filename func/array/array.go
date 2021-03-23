package array

import (
	"encoding/json"
)

func Keys(cmap interface{}) []string {

	var keys []string
	content, err := json.Marshal(cmap)
	if err == nil {
		maps := make(map[string]interface{})
		json.Unmarshal(content, &maps)
		for key, _ := range maps {
			keys = append(keys, key)
		}
	}
	return keys
}

func Search(obj interface{}, value interface{}) {

}
