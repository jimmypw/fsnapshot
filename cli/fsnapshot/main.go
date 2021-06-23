package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jimmypw/fsnapshot"
)

var dir *string
var out *string
var compare *string

func initflags() {
	dir = flag.String("dir", "", "The directory containing files you want to snapshot")
	out = flag.String("out", "", "Save the new manifest to this location")
	compare = flag.String("compare", "", "The existing manifest you would like to compare the directory against")
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	initflags()
	flag.Parse()

	if !isFlagPassed("dir") {
		fmt.Println("You must provide a directory to scan")
	}

	if !isFlagPassed("out") && !isFlagPassed("compare") {
		fmt.Println("You must pass either one of -out or -compare")
	}

	fm, err := fsnapshot.Snapshot(*dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if isFlagPassed("out") {
		fmt.Println("Saving")
		fm.Save(*out)
	}

	if isFlagPassed("compare") {
		var fm2 fsnapshot.FileManifest
		fm2.Load(*compare)
		fmreport, err := fsnapshot.Compare(&fm2, &fm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmreport.Report()
	}

}
