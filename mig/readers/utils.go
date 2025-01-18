package readers

import (
	"go.reizu.org/greek"
	"strings"
)

func ToGreeklish(name string) string {
	name = strings.TrimSpace(name)
	code := greek.Greeklish(strings.ToLower(name))
	code = nonAlphanumericRegex.ReplaceAllString(code, "")
	code = strings.Replace(code, " ", "-", -1)
	code = strings.Replace(code, "&", "", -1)
	return code
}
