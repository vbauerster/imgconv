// Exercise 10.1: Extend the jpeg program so that it converts any supported
// input format to any output format, using image.Decode to detect the input
// format and a flag to select the output format.

package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Format string `short:"f" long:"format" description:"output format" value-name:"png|jpg"`
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS] file1 file2 ..."
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
	for _, infile := range args {
		outfile, err := convert(infile, opts.Format)
		verboseMsg := infile + " => "
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s%v\n", verboseMsg, err)
			continue
		}
		fmt.Printf("%s%s\n", verboseMsg, outfile)
	}
}

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

func convert(infile, format string) (string, error) {
	format = strings.ToLower(format)
	outfile := strings.TrimSuffix(infile, filepath.Ext(infile)) + "." + format
	return outfile, convert2(outfile, infile, format)
}

func convert2(outfile, infile, format string) error {
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
