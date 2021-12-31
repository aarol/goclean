package util

import "fmt"

func HumanizeBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	unit := 0
	for bytes >= 1024 {
		bytes /= 1024
		unit++
	}
	return fmt.Sprintf("%d%s", bytes, units[unit])
}
