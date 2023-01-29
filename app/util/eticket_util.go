package util

import (
	"fmt"
	"github.com/frchandra/gmcgo/config"
	"github.com/h2non/bimg"
	"github.com/skip2/go-qrcode"
)

type ETicketUtil struct {
	config *config.AppConfig
}

func NewETicketUtil(config *config.AppConfig) *ETicketUtil {
	return &ETicketUtil{config: config}
}

func (e *ETicketUtil) GenerateETicket(seatName, seatLink string) error {

	url := e.config.AppUrl + ":" + e.config.AppPort + "/api/v1/seat/" + seatLink
	qr, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	title := bimg.Watermark{
		Top:         2,
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
		Top:         13,
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
		return err
	}

	ticket, err := bimg.NewImage(frame).WatermarkImage(content)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := bimg.Write("./storage/ticket/"+seatName+".png", ticket); err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("./storage/ticket/" + seatName + ".png GENERATED===")

	return nil
}
