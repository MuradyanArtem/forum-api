package infrastructure

import (
	"io"

	"github.com/sirupsen/logrus"
)

func ConfigureLogger(writer io.Writer, l string) error {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})

	logrus.SetOutput(writer)

	level, err := logrus.ParseLevel(l)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	return nil
}
