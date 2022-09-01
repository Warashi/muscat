//go:build golangdesign

package clipboard

import "golang.design/x/clipboard"

func init() {
	if err := clipboard.Init(); err != nil {
		unsupported = true
	}
}

var (
	unsupported bool
)

func Unsupported() bool {
	return unsupported
}

func Write(data []byte) error {
	if Unsupported() {
		return ErrUnsupported
	}
	<-clipboard.Write(clipboard.FmtText, data)
	return nil
}

func Read() ([]byte, error) {
	if Unsupported() {
		return nil, ErrUnsupported
	}
	return clipboard.Read(clipboard.FmtText), nil
}
