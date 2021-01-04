package main

import (
	"github.com/qpliu/qrencode-go/qrencode"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-colorable"
	"io/ioutil"
	"os"

	"github.com/vandr0iy/qrc/lib"
	"github.com/vandr0iy/go-tty"
)

type cmdOptions struct {
	Help    bool `short:"h" long:"help" description:"show this help message"`
	Inverse bool `short:"i" long:"invert" description:"invert color"`
}

func showHelp() {
	const v = `Usage: qrc [OPTIONS] [TEXT]

Options:
  -h, --help
    Show this help message
  -i, --invert
    Invert color

Text examples:
  http://www.example.jp/
  MAILTO:foobar@example.jp
  WIFI:S:myssid;T:WPA;P:pass123;;
`

	os.Stderr.Write([]byte(v))
}

func pErr(format string, a ...interface{}) {
	fmt.Fprint(os.Stdout, os.Args[0], ": ")
	fmt.Fprintf(os.Stdout, format, a...)
}

func main() {
	ret := 0
	defer func() { os.Exit(ret) }()

	opts := &cmdOptions{}
	optsParser := flags.NewParser(opts, flags.PrintErrors)
	args, err := optsParser.Parse()
	if err != nil || len(args) > 1 {
		showHelp()
		ret = 1
		return
	}
	if opts.Help {
		showHelp()
		return
	}

	var text string
	if len(args) == 1 {
		text = args[0]
	} else {
		text_bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			pErr("read from stdin failed: %v\n", err)
			ret = 1
			return
		}
		text = string(text_bytes)
	}

	grid, err := qrencode.Encode(text, qrencode.ECLevelL)
	if err != nil {
		pErr("encode failed: %v\n", err)
		ret = 1
		return
	}

	da1, err := tty.GetDeviceAttributes1(os.Stdout)
	if err == nil && da1[tty.DA1_SIXEL] {
		qrc.PrintSixel(os.Stdout, grid, opts.Inverse)
	} else {
		stdout := colorable.NewColorableStdout()
		qrc.PrintAA(stdout, grid, opts.Inverse)
	}
}
