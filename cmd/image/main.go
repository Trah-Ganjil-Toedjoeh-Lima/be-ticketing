package main

import (
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/h2non/bimg"
	"github.com/skip2/go-qrcode"
)

func main() {
	seatName := "H33"
	url := config.NewAppConfig().AppUrl + config.NewAppConfig().AppPort
	qr, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		fmt.Println(err.Error())

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
		Text:        "SEAT - " + seatName,
		Opacity:     1,
		Width:       200,
		DPI:         100,
		Margin:      0,
		Font:        "sans bold 12",
		Background:  bimg.Color{0, 0, 0},
		NoReplicate: true,
	}

	newQrCode, err := bimg.NewImage(qr).Watermark(title)
	newQrCode, err = bimg.NewImage(newQrCode).Watermark(seat)

	content := bimg.WatermarkImage{
		Left:    0,
		Top:     0,
		Buf:     newQrCode,
		Opacity: 1,
	}

	frame, err := bimg.Read("./storage/picture/polite_cat.png")
	if err != nil {
		fmt.Println(err.Error())

	}

	ticket, err := bimg.NewImage(frame).WatermarkImage(content)
	if err != nil {
		fmt.Println(err.Error())

	}

	if err := bimg.Write("./storage/ticket/"+seatName+".png", ticket); err != nil {
		fmt.Println(err.Error())

	}

}
