package util

import "github.com/wobeproject/logger"

//ReverseString ...
func ReverseString(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

//FailOnError ...
func FailOnError(msg string, err error) {
	if err != nil {
		l := logger.GetInstance()
		l.Error(msg, map[string]interface{}{
			"ERR": err.Error(),
		})
		//		panic(err)
	}

}
