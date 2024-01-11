package funcs

import (
	"gopkg.in/yaml.v3"
	"strings"
)

func ToYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

func Indent(indent, s string) string {
	if indent == "" || s == "" {
		return s
	}
	lines := strings.SplitAfter(s, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(append([]string{""}, lines...), indent)
}
