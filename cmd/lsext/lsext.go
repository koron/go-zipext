package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"log"

	"github.com/koron/go-zipext"
)

func lsext(name string) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, zf := range zr.File {
		fmt.Printf("%s\n", zf.Name)
		ext := zipext.NewReader(zf)
		fmt.Println("  Extra fields:n")
		for {
			f, err := ext.Read()
			if err != nil {
				return err
			}
			if f == nil {
				break
			}
			fmt.Printf("    %04x: %-32x\n", f.Tag, f.Data)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	for _, s := range flag.Args() {
		if err := lsext(s); err != nil {
			log.Fatal(err)
		}
	}
}
