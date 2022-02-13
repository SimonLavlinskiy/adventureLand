package helpers

import "strconv"

func ToString(int int) string {
	return strconv.FormatInt(int64(int), 10)
}

func ToInt(string string) int {
	numInt64, _ := strconv.ParseUint(string, 10, 32)
	return int(numInt64)
}
