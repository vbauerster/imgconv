// Exercise 10.1: Extend the jpeg program so that it converts any supported
// input format to any output format, using image.Decode to detect the input
// format and a flag to select the output format.

package main

import (
	"fmt"
	"gopl/ch10/ex101/imgconv"
	"os"
	"sync"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Format  string `short:"f" long:"format" description:"output format" value-name:"png|jpg"`
	Verbose bool   `short:"v" long:"verbose" description:"Verbose progress messages"`
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS] file1 file2 ..."
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	results := make(chan result)
	var wg sync.WaitGroup
	for _, filename := range args {
		wg.Add(1)
		go convertImage(filename, opts.Format, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if opts.Verbose {
			verboseMsg := res.infile + " => "
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s%v\n", verboseMsg, res.err)
				continue
			}
			fmt.Printf("%s%s\n", verboseMsg, res.outfile)
		}
	}

}

func incoming(args []string) <-chan string {
	c := make(chan string)
	go func() {
		for _, filename := range args {
			c <- filename
		}
		close(c)
	}()
	return c
}

type result struct {
	infile  string
	outfile string
	err     error
}

func convert(filenames []string, format string, results chan<- result) {
	var wg sync.WaitGroup
	for _, f := range filenames {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			var r result
			r.outfile, r.err = imgconv.ConvertToSameDir(f, format)
			results <- r
		}(f)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
}

var sema = make(chan struct{}, 20)

func convertImage(filename, format string, wg *sync.WaitGroup, results chan<- result) {
	sema <- struct{}{}
	defer func() {
		<-sema
		wg.Done()
	}()
	r := result{infile: filename}
	r.outfile, r.err = imgconv.ConvertToSameDir(filename, format)
	results <- r
}
