package collector

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

const (
	fileName = "ngraph.prom"
)

type FileWriter struct {
	OutputFile *os.File
	Writer     *bufio.Writer
}

func NewFileWriter(path string) (*FileWriter, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	output, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(output)

	return &FileWriter{
		OutputFile: output,
		Writer:     writer,
	}, nil

}

func (fw *FileWriter) Write(metric string) error {
	_, err := fw.Writer.WriteString(metric)
	if err != nil {
		return err
	}

	return nil
}

func (fw *FileWriter) Close() error {
	if err := fw.Writer.Flush(); err != nil {
		return err
	}

	return fw.OutputFile.Close()
}

// writeToFile writes network connection metrics to a file specified by the given path.
// It takes a map of connection groups and their counts, and formats and writes Prometheus metrics to the specified file.
// The function returns an error if there is an issue creating or writing to the file.
func (c *Collector) writeToFile(path string, connections map[string]int) error {
	writer, err := NewFileWriter(path)
	if err != nil {
		return err
	}
	defer writer.Close()

	for k, v := range connections {
		metric, err := createMetric(k, v)
		if err != nil {
			return fmt.Errorf("error when create metric(%s): %w", k, err)
		}

		if err := writer.Write(metric); err != nil {
			return fmt.Errorf("error when writes a metric(%s): %w", metric, err)
		}
	}

	return nil
}
