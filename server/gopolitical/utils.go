package gopolitical

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func PositiveModulus(a, b int) int {
	return (a%b + b) % b
}

func Info(name string, format string, args ...any) {
	args = append([]any{name}, args...)
	format = strings.ReplaceAll(format, "%s", "\033[1;33m%s\033[0m\033[30m")
	format = strings.ReplaceAll(format, "%-24s", "\033[1;33m%-24s\033[0m\033[30m")
	format = strings.ReplaceAll(format, "%d", "\033[1m\033[35m%d\033[0m\033[30m")
	format = strings.ReplaceAll(format, "%f", "\033[1m\033[35m%.2f\033[0m\033[30m")
	format = strings.ReplaceAll(format, "%.0f", "\033[1m\033[35m%.0f\033[0m\033[30m")
	format = strings.ReplaceAll(format, "%.2f", "\033[1m\033[35m%.2f\033[0m\033[30m")
	format = strings.ReplaceAll(format, "%v", "\033[32m%v\033[0m\033[30m")
	format = "[\033[36m%-13s\033[0m] \033[30m" + format + "\033[0m"
	log.Printf(format, args...)
}

func Debug(name string, format string, args ...any) {
	if DebugEnabled() {
		Info(name, format, args...)
	}
}

func DebugEnabled() bool {
	envVariable := os.Getenv("GOPOLITICAL_DEBUG")
	return envVariable != "" && envVariable != "0"
}

// Global random generator
var random = rand.New(rand.NewSource(time.Now().UnixNano()))
