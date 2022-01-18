package bot

import (
	"bytes"
	"image"
	"image/color"
	"net/http"

	"github.com/AidenHadisi/MyDailyBibleBot/types"
	"github.com/fogleman/gg"
)

type ImageProcessor struct {
	client types.HttpClient
}

func NewImageProcessor(client types.HttpClient) *ImageProcessor {
	return &ImageProcessor{
		client: client,
	}
}

func (p *ImageProcessor) Process(url, text string) ([]byte, error) {
	bgImage, err := p.getImage(url)
	if err != nil {
		return nil, err
	}

	//first draw the image
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)

	//then draw a semi transparent rectangle over the image to make it darker
	dc.SetColor(color.RGBA{0, 0, 0, 100})
	dc.DrawRectangle(0, 0, float64(imgWidth), float64(imgHeight))
	dc.Fill()

	//Now draw the text inside rectangle
	x := float64(imgWidth / 2)
	y := float64((imgHeight / 2) - 80)
	maxWidth := float64(imgWidth) - 60.0
	dc.SetColor(color.White)
	dc.LoadFontFace("KeepCalm-Medium.ttf", float64(30))
	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	var buf bytes.Buffer
	err = dc.EncodePNG(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *ImageProcessor) getImage(url string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}
