package plugin

import (
	"html/template"
	"strconv"
)

// Arithmetic returns a template.FuncMap
func Arithmetic() template.FuncMap {
	f := make(template.FuncMap)

	f["TIMES"] = func(numberA, numberB float64) string {
		return strconv.FormatFloat(numberA*numberB, 'f', 1, 64)
	}
	f["MINUS"] = func(numberA, numberB int) string {
		return strconv.Itoa(numberA - numberB)
	}

	return f
}
