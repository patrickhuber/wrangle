package structio

import "io"

// JSONOutput is the format output for json
const JSONOutput = "json"

// YAMLOutput is the format output for yaml
const YAMLOutput = "yaml"

// TableOutput is the format output for table
const TableOutput = "table"

// NewWriter returns a new writer for the given output type
func NewWriter(writer io.Writer, output string) Writer {
	switch output {
	case JSONOutput:
		return NewJSONWriter(writer)
	case YAMLOutput:
		return NewYamlWriter(writer)
	default:
		return NewTableWriter(writer)
	}
}
