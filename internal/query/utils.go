package query

import "strings"

func parseSort(sortString string) (string, int) {
	pair := strings.Split(sortString, ":")
	if len(pair) != 2 {
		return "", 0
	}

	key := pair[0]
	order := pair[1]
	value := 0

	if order == "asc" {
		value = 1
	} else if order == "desc" {
		value = -1
	} else {
		return "", 0
	}

	return key, value
}
