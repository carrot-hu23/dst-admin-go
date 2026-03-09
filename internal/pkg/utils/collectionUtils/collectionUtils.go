package collectionUtils

import "log"

func ToSet(list []string) []string {

	var m = map[string]string{}
	var set []string

	for i, _ := range list {
		key := list[i]
		log.Println("key", key)
		_, ok := m[key]
		if !ok {
			m[key] = key
			set = append(set, key)
		}
	}
	return set
}
