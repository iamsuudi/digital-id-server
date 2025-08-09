package utils

import (
	"fmt"
	"path/filepath"
	"time"
)

// MakeFileName generates a unix time filename for a file passed to it
func MakeFileName(original string) string {
	ext := filepath.Ext(original)            // ".jpg", ".png", etc.
	seconds := time.Now().UnixNano()             // e.g. 1717923845
	return fmt.Sprintf("%d%s", seconds, ext) // 1717923845.jpg
}
