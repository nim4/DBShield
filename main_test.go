package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

func Test(t *testing.T) {
	usage(true)

	//
	os.Args = []string{os.Args[0], "-k", "-c", "conf/dbshield.yml"}
	main()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{os.Args[0], "-k", "-c", "conf/invalid.yml"}
	main()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{os.Args[0], "-l"}
	main()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{os.Args[0], "-version"}
	main()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{os.Args[0], "-h"}
	main()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	dat, err := ioutil.ReadFile("conf/dbshield.yml")
	if err != nil {
		t.Fatal(err)
	}
	dat = bytes.Replace(dat, []byte("dbDir: "), []byte("dbDir: /INVALID"), 1)
	path := os.TempDir() + "/tempconfig.yml"
	err = ioutil.WriteFile(path, dat, 0600)
	if err != nil {
		t.Fatal(err)
	}
	os.Args = []string{os.Args[0], "-k", "-c", path}
	main()
}
