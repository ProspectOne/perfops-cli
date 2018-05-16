// +build !windows

package internal

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"

	"github.com/gosuri/uilive"
)

func termSize() (int, int) {
	if termCols == 0 {
		if out, err := os.OpenFile("/dev/tty", syscall.O_WRONLY, 0); err == nil {
			defer out.Close()
			var size struct {
				rows    uint16
				cols    uint16
				xpixels uint16
				ypixels uint16
			}
			_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, out.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&size)))
			termCols, termRows = int(size.cols), int(size.rows)
		}
	}
	return termCols, termRows
}

func termStartOfRow(w io.Writer) {
	w.Write([]byte(fmt.Sprintf("%c[%dD", uilive.ESC, termCols)))
}
