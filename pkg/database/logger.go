package database

import (
	"fmt"
	"io"
)

type (
	loggerWriter struct{ io.Writer }
)

func (l loggerWriter) Printf(format string, data ...any) {
	fmt.Fprintf(l.Writer, format, data...)
}
