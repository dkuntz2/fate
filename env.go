package fate

import (
	"fmt"
	"os"
)

func EnvValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Errorf("No `%s` environment variable set", key))
	}

	return value
}
