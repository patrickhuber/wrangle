package enums

import (
	"fmt"
	"strings"
)

type EnumValue struct {
	Enum     []string
	Default  string
	selected string
}

func (e *EnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			e.selected = value
			return nil
		}
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e EnumValue) String() string {
	if e.selected == "" {
		return e.Default
	}
	return e.selected
}

func NewFormatEnum() *EnumValue {
	return &EnumValue{
		Enum:    []string{"json", "table", "yaml"},
		Default: "table",
	}
}
