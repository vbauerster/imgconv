// Package imgconv provides image conversion facility
package imgconv

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ConvertImg(in io.Reader, out io.Writer, format string, quality int) error {
	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}
	switch format {
	case "jpg":
		if quality < 0 {
			quality = 90
		}
		return jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
	case "png":
		return png.Encode(out, img)
	case "gif":
		return gif.Encode(out, img, &gif.Options{NumColors: quality})
	default:
		return newErrorf(ErrUnsupportedFormat, "unsupported format: %q", format)
	}
}

// ConvertToSameDir reads an image from infile and writes
// a converted version of it in the same directory.
// It returns the generated file name, e.g. "foo.jpg".
func ConvertToSameDir(infile, format string, quality int) (string, error) {
	format = strings.ToLower(format)
	inExt := filepath.Ext(infile)
	outExt := "." + format
	basename := strings.TrimSuffix(infile, inExt)
	var outfile string
	if strings.ToLower(inExt) != outExt {
		outfile = basename + outExt
	} else {
		outfile = basename + "_conv" + outExt
	}
	return outfile, Convert(infile, outfile, format, quality)
}

func Convert(infile, outfile, format string, quality int) error {
	in, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outfile)
	if err != nil {
		return err
	}

	err = ConvertImg(in, out, format, quality)
	if closeErr := out.Close(); err == nil {
		err = closeErr
	} else if err, ok := err.(*Error); ok && err.Type == ErrUnsupportedFormat {
		os.Remove(outfile)
	}

	return err
}
