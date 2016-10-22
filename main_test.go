package main

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	usage(true)

	os.Args = []string{os.Args[0], "-k", "-c", "conf/dbshield.yml"}
	main()
}
