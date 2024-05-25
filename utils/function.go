package utils

import "fmt"

/* ========== Resource Function ========== */

func IsLabelEqual(a map[string]string, b map[string]string) bool {
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

/* ========== Time Function ========== */

func WaitForever() {
	<-make(chan struct{})
}

/* ========== Rand Function ========== */

func GenerateName(name string, n int) string {
	return fmt.Sprint(name, "-", n)
}
