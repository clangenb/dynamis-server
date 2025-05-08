package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	enc := flag.Bool("encrypt", false, "Encrypt input file")
	dec := flag.Bool("decrypt", false, "Decrypt input file")
	in := flag.String("in", "", "Input file path")
	out := flag.String("out", "", "Output file path")

	flag.Parse()

	if *in == "" || *out == "" {
		fmt.Println("Usage: -encrypt/-decrypt -in input -out output")
		os.Exit(1)
	}

	var err error
	if *enc {
		err = EncryptFile(*in, *out)
	} else if *dec {
		err = DecryptFile(*in, *out)
	} else {
		fmt.Println("Specify -encrypt or -decrypt")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Success!")
}
