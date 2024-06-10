package router

import "mime"

var ext2type = map[string]string{
	".js":   "text/javascript; charset=utf-8",
	".css":  "text/css; charset=utf-8",
	".html": "text/html; charset=utf-8",
}

// TypeByExtension
// On my machine mime.TypeByExtension for some reason returns 'text/plain; charset=utf-8'
func TypeByExtension(extension string) string {
	if v, h := ext2type[extension]; h {
		return v
	} else {
		return mime.TypeByExtension(extension)
	}
}
