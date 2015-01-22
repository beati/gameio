package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"log"
	"os"
	"strings"
)

var (
	packageName string
	outFileName string
)

func main() {
	flag.StringVar(&packageName, "p", "assets", "pacakge name")
	flag.StringVar(&outFileName, "o", "", "output file")
	flag.Parse()
	if outFileName == "" {
		outFileName = packageName + ".go"
	}

	if len(flag.Args()) == 0 {
		log.Fatal("no input files")
	}

	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	out := bufio.NewWriter(outFile)

	_, err = fmt.Fprintf(out, "package %s\n\n", packageName)
	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()
	for i, inFileName := range args {
		if !strings.HasSuffix(inFileName, ".png") {
			log.Println(inFileName, ": not a png file")
			continue
		}

		inFile, err := os.Open(inFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer inFile.Close()
		in := bufio.NewReader(inFile)

		img, _, err := image.Decode(in)
		if err != nil {
			log.Println(inFileName, ": bad image format")
			continue
		}

		err = convert(img, stripName(inFileName), out)
		if err != nil {
			log.Fatal(err)
		}

		if i != len(args)-1 {
			_, err = fmt.Fprintf(out, "\n\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = out.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

const structHeaderFormat = `var %s = struct {
	W int
	H int
	Pixels [%d]uint32
}{%d, %d, [%d]uint32{
`

func convert(img image.Image, name string, out io.Writer) error {
	b := img.Bounds()
	w := b.Max.X - b.Min.X
	h := b.Max.Y - b.Min.Y
	_, err := fmt.Fprintf(out, structHeaderFormat, name, w*h, w, h, w*h)
	if err != nil {
		return err
	}

	pixelCount := 0
	const columnCount = 7
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			const scale = 0xFFFF / 0xFF
			r /= scale
			g /= scale
			b /= scale
			a /= scale
			c := r<<24 | g<<16 | b<<8 | a
			_, err = fmt.Fprintf(out, "%#08x, ", c)
			if err != nil {
				return err
			}
			pixelCount++
			if pixelCount%columnCount == 0 {
				_, err = fmt.Fprintf(out, "\n")
				if err != nil {
					return err
				}
			}
		}
	}

	if pixelCount%columnCount != 0 {
		_, err = fmt.Fprintf(out, "\n")
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(out, "}}")
	if err != nil {
		return err
	}

	return nil
}

func stripName(fileName string) string {
	return strings.Title(strings.TrimSuffix(fileName, ".png"))
}
