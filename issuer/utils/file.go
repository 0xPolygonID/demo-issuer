package utils

import (
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func ReadFileByPath(basePath string, fileName string) ([]byte, error) {
	logger.Debug("utils.ReadFileByPath() invoked")

	path := filepath.Join(basePath, fileName)
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrapf(err, "failed open file '%s' by path '%s'", fileName, path)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed read file '%s' by path '%s'", fileName, path)
	}
	return data, nil
}
