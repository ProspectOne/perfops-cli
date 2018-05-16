package internal

import (
	"io"
)

func termSize() (int, int) {
	// TODO: Implement me.
	return termCols, termRows
}

func termStartOfRow(w io.Writer) {
	// TODO: Implement me.
}
