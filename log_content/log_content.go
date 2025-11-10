package log_content

import "fmt"

func LogContentNormal(tag string, message string, args ...any) string {
	return fmt.Sprintf("[%s]%s", tag, fmt.Sprintf(message, args...))
}
