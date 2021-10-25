package structio

import (
	"encoding/json"
	"io"
)

type jsonWriter struct {
	writer io.Writer
}

// NewJSONWriter returns a Writer that outputs json
func NewJSONWriter(writer io.Writer) Writer {
	return &jsonWriter{
		writer: writer,
	}
}

func (w *jsonWriter) Write(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(bytes)
	return err
}
