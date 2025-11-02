package main

import (
	"os"
	"strings"
)

func allowedOrigins() []string {
	allowed := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowed) == 1 && allowed[0] == "" {
		allowed = []string{"*"}
	}
	return allowed
}
