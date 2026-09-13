package procesamiento

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const maxBytesPermitidos = 20 * 1024 * 1024 // 20 MB

func Validar(path string) (mime string, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.Size() > maxBytesPermitidos {
		return "", fmt.Errorf("supera el tamaño máximo permitido (%d MB)", maxBytesPermitidos/1024/1024)
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	return http.DetectContentType(buf[:n]), nil
}
