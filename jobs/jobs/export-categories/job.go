package export_categories

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
)

const exportFile = "program-data/export_categories.json"

type job struct {
}

func newJob() *job {
	return &job{}
}

func (j *job) Run() {
	var l = slog.With("module", "jobs/jobs/export-categories", "method", "Run")
	// Check if the export file exists:
	if _, e := os.Stat(exportFile); e == nil {
		return
	}

	// Create the base export file:
	_, e := os.Create(filepath.Dir(exportFile))
	data, e := newExporter().Export()
	if e != nil {
		l.Warn("failed to export categories", "error", e)
		return
	}
	// Write the data to the export file:
	jsonData, e := prettyMarshal(data)
	if e != nil {
		l.Warn("failed to marshal data", "error", e)
		return
	}
	e = os.WriteFile(exportFile, jsonData, 0600)
	if e != nil {
		l.Warn("failed to write data", "error", e)
	}
	l.Info("exported categories to file", "file", exportFile)
}

func prettyMarshal(v any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	e := encoder.Encode(v)
	return buffer.Bytes(), e
}
