package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"os"

	"github.com/nfnt/resize"
	qrcode "github.com/skip2/go-qrcode"
)

func generateQRWithLogo(url string) image.Image {
	qrCode, err := qrcode.New(url, qrcode.Highest)
	if err != nil {
		fmt.Println("Failed to generate QR code:", err)
		return nil
	}
	qrImg := qrCode.Image(512)

	logoFile, err := os.Open("logo.png")
	if err != nil {
		fmt.Println("Failed to open logo:", err)
		return nil
	}
	defer logoFile.Close()

	logoImg, _, err := image.Decode(logoFile)
	if err != nil {
		fmt.Println("Failed to decode logo:", err)
		return nil
	}

	logoSize := 75
	logoImg = resize.Resize(uint(logoSize), uint(logoSize), logoImg, resize.Lanczos3)

	offset := image.Pt((qrImg.Bounds().Dx()-logoSize)/2, (qrImg.Bounds().Dy()-logoSize)/2)
	m := image.NewRGBA(qrImg.Bounds())
	draw.Draw(m, m.Bounds(), qrImg, image.Point{}, draw.Over)
	draw.Draw(m, logoImg.Bounds().Add(offset), logoImg, image.Point{}, draw.Over)

	return m
}

func qrHandler(w http.ResponseWriter, r *http.Request) {
	urls, ok := r.URL.Query()["url"]
	if !ok || len(urls[0]) < 1 {
		http.Error(w, "Url parameter is missing", http.StatusBadRequest)
		return
	}
	url := urls[0]

	img := generateQRWithLogo(url)
	if img == nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)
}

func main() {
	http.HandleFunc("/qr", qrHandler)
	fmt.Println("Server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
