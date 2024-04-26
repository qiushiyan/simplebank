package report

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/qiushiyan/go-quarto"
	"github.com/qiushiyan/simplebank/foundation/random"
)

type QuartoGenerator struct{}

func NewQuartoGenerator() *QuartoGenerator {
	return &QuartoGenerator{}
}

// New generates a new quarto report based on subject and data, returns the output file name and error
func (g QuartoGenerator) New(ctx context.Context, subject string, data any) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	outputFileName := os.TempDir() + "report-" + random.RandomString(8) + ".pdf"
	quartoConfig := &quarto.Config{
		Output: outputFileName,
	}

	var prefix string
	if testing.Testing() {
		prefix = filepath.Join(cwd, "..", "..", "zarf", "quarto", "templates")
	} else {
		prefix = filepath.Join(cwd, "zarf", "quarto", "templates")
	}

	var templateFilePath string
	switch subject {
	case SubjectActivity:
		templateFilePath = filepath.Join(prefix, "activity.qmd")
		d, ok := data.(SubjectActivityData)
		if !ok {
			return "", errors.New("invalid type for Data property in SubjectActivity")
		}
		quartoConfig.SetParam("user_name", d.Username)
	default:
		return "", ErrInvalidSubject
	}

	_, err = quarto.Render(ctx, templateFilePath, quartoConfig)
	if err != nil {
		return "", err
	}

	return outputFileName, nil
}
