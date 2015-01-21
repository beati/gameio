package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"log"
	"os"
)

func main() {
	infile, err := os.Open("test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()
	in := bufio.NewReader(infile)

	outfile, err := os.Create("image.g")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	out := bufio.NewWriter(outfile)

	err = convert(in, "img", out)
	if err != nil {
		log.Fatal(err)
	}

	err = out.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func convert(in io.Reader, name string, out io.Writer) error {
	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "var %s = []uint32{\n", name)
	if err != nil {
		return err
	}

	b := img.Bounds()
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

	_, err = fmt.Fprintf(out, "}")
	if err != nil {
		return err
	}

	return nil
}
