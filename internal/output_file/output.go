package outputfile

import (
	"log/slog"
	"os"
	"strings"
)

type OutputFile struct {
	fileName string
	content  strings.Builder
}

func NewOutputFile(fileName string) (*OutputFile, error) {
	if len(fileName) > 0 {
		err := os.WriteFile(fileName, []byte{}, 0644)
		if err != nil {
			slog.Error("Failed creating outputfile", slog.String("file", fileName), slog.Any("error", err))
			return nil, err
		}
		os.Remove(fileName)
		slog.Debug("Output", slog.String("file", fileName))
	} else {
		slog.Debug("Output disabled")
	}
	return &OutputFile{
		fileName: fileName,
		content:  strings.Builder{},
	}, nil
}

func (o *OutputFile) SetFile(fileName string) error {
	if len(fileName) > 0 {
		err := os.WriteFile(fileName, []byte{}, 0644)
		if err != nil {
			slog.Error("Failed creating outputfile", slog.String("file", fileName), slog.Any("error", err))
			return err
		}
		os.Remove(fileName)
		o.fileName = fileName
		o.content.Reset()
		slog.Debug("Output", slog.String("file", fileName))
	} else {
		slog.Debug("Output disabled")
	}
	return nil
}

func (o *OutputFile) AddContent(line string) {
	o.content.WriteString(line)
}

func (o *OutputFile) Close() error {
	if len(o.fileName) == 0 {
		slog.Debug("Output closed [empty]")
		return nil
	}
	err := os.WriteFile(o.fileName, []byte(o.content.String()), 0644)
	if err != nil {
		slog.Error("Output save error", slog.String("file", o.fileName), slog.Any("error", err), slog.String("content", o.content.String()))
	} else {
		slog.Debug("Output save", slog.String("file", o.fileName), slog.String("content", o.content.String()))
	}
	return err
}
