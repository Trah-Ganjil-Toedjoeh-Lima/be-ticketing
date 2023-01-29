package main

import (
	"fmt"
	"github.com/h2non/bimg"
	"github.com/skip2/go-qrcode"
	"os"
)

func main() {
	_ = qrcode.WriteFile("https://example.org", qrcode.Medium, 256, "qr_example.png")

	qrcode, err := bimg.Read("./storage/picture/qr_example.png")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	title := bimg.Watermark{
		Top:         3,
		Left:        38,
		Text:        "GRAND CONCERT GMCO 2023",
		Opacity:     1,
		Width:       200,
		DPI:         100,
		Margin:      0,
		Font:        "sans bold 9",
		Background:  bimg.Color{0, 0, 0},
		NoReplicate: true,
	}

	seat := bimg.Watermark{
		Top:         15,
		Left:        90,
		Text:        "SEAT - H33",
		Opacity:     1,
		Width:       200,
		DPI:         100,
		Margin:      0,
		Font:        "sans bold 12",
		Background:  bimg.Color{0, 0, 0},
		NoReplicate: true,
	}

	newQrCode, err := bimg.NewImage(qrcode).Watermark(title)
	newQrCode, err = bimg.NewImage(newQrCode).Watermark(seat)

	content := bimg.WatermarkImage{
		Left:    0,
		Top:     0,
		Buf:     newQrCode,
		Opacity: 1,
	}

	frame, err := bimg.Read("./storage/picture/polite_cat.png")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	ticket, err := bimg.NewImage(frame).WatermarkImage(content)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	bimg.Write("./storage/ticket/new_ticket.png", ticket)

}
