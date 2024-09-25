package decorator

import (
	"fmt"
)

const (
	KeyprefixToken = "thread:token"
)

func Token(key string) string {
	return fmt.Sprintf("%s:%s", KeyprefixToken, key)
}
