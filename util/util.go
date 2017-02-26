package util

//ReverseString ...
func ReverseString(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
//	result = result + " new"
	return
}
