package handler

import (
	"fmt"
)

func httpData(data string) string {
	return fmt.Sprintf("{\"data\": %s}", data)
}
