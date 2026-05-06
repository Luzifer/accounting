package database

import (
	"fmt"
	"io"
)

type (
	loggerWriter struct{ io.Writer }
)

func (l loggerWriter) Printf(format string, data ...any) {
	_, _ = fmt.Fprintf(l.Writer, format, data...)
}
