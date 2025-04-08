package util

import (
	"os"
	"strings"
)

func GetRealValue(v string) string {
	if strings.HasPrefix(v, "env:") {
		return os.Getenv(v[4:])
	}

	return v
}
