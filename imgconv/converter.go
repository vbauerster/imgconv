// Package imgconv provides image conversion facility
package imgconv

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func convertImg(in io.Reader, out io.Writer, format string) error {
	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}
	// fmt.Fprintln(os.Stderr, "Input format =", kind)
	switch format {
	case "jpg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "png":
		return png.Encode(out, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// ConvertToSameDir reads an image from infile and writes
// a converted version of it in the same directory.
// It returns the generated file name, e.g. "foo.jpg".
func ConvertToSameDir(infile, format string) (string, error) {
	format = strings.ToLower(format)
	outfile := strings.TrimSuffix(infile, filepath.Ext(infile)) + "." + format
	return outfile, Convert(infile, outfile, format)
}

func Convert(infile, outfile, format string) error {
	in, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outfile)
	if err != nil {
		return err
	}
	if err := convertImg(in, out, format); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
