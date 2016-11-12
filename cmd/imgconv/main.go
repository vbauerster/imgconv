// Exercise 10.1: Extend the jpeg program so that it converts any supported
// input format to any output format, using image.Decode to detect the input
// format and a flag to select the output format.

package main

import (
	"fmt"
	"gopl/ch10/ex101/imgconv"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/jessevdk/go-flags"
)

const (
	_  = iota
	kb = 1 << (iota * 10)
)

type options struct {
	Format  string `short:"o" long:"outfmt" description:"Output format. If format supports quality setting, you may specify one after colon like jpg:90" value-name:"png|jpg" required:"true"`
	Verbose bool   `short:"v" long:"verbose" description:"Verbose progress messages"`
}

type result struct {
	infile  string
	outfile string
	delta   int64
	err     error
}

var sema = make(chan struct{}, 20)

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "imgconv"
	parser.Usage = "[OPTIONS] file1 file2 ..."
	args, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			parser.WriteHelp(os.Stderr)
			os.Exit(1)
		}
	}

	if len(args) == 0 {
		fmt.Printf("nothing to do\n\n")
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	format := opts.Format
	quality := 90
	colonIndex := strings.Index(opts.Format, ":")
	if colonIndex > 0 {
		format = opts.Format[:colonIndex]
		qstr := opts.Format[colonIndex+1:]
		if len(qstr) > 0 {
			quality, err = strconv.Atoi(qstr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid qualtiy value: %q", qstr)
				os.Exit(1)
			}
		}
	}

	if len(strings.Trim(format, " ")) == 0 {
		fmt.Printf("unexpected argument for flag '-o, --outfmt'\n\n")
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	results := make(chan result)
	var wg sync.WaitGroup
	for _, filename := range args {
		wg.Add(1)
		go convert(filename, format, quality, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if opts.Verbose {
			verboseMsg := res.infile + " => "
			if res.err != nil {
				fmt.Fprintf(os.Stderr, "%s%v\n", verboseMsg, res.err)
				continue
			}
			if res.delta > 0 {
				fmt.Printf("%s%s (save: %.1f Kb)\n", verboseMsg, res.outfile, float64(res.delta)/kb)
			} else {
				fmt.Printf("%s%s\n", verboseMsg, res.outfile)
			}
		}
	}
}

func convert(filename, format string, quality int, wg *sync.WaitGroup, results chan<- result) {
	sema <- struct{}{}
	defer func() {
		<-sema
		wg.Done()
	}()
	r := result{infile: filename}
	r.outfile, r.err = imgconv.ConvertToSameDir(filename, format, quality)
	if r.err == nil {
		in, _ := os.Stat(r.infile)
		out, _ := os.Stat(r.outfile)
		r.delta = in.Size() - out.Size()
	}
	results <- r
}
