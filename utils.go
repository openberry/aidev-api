package aidev

import (
	"bytes"
	"image/png"

	"github.com/boombuler/barcode/qr"
)

// Generates a new QR code as a PNG image
func QRCode(data string) ([]byte, error) {
	// Generate code
	code, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	// Encode the barcode as png
	buf := bytes.NewBuffer(nil)
	err = png.Encode(buf, code)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
