package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"bytes"
)

// ListLocations -
func (c *Client) ListLocations(pageURL *string) (RespShallowLocations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if val, ok := c.cache.Get(url); ok {
		locationsResp := RespShallowLocations{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			return RespShallowLocations{}, err
		}

		return locationsResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespShallowLocations{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespShallowLocations{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespShallowLocations{}, err
	}

	locationsResp := RespShallowLocations{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return RespShallowLocations{}, err
	}

	c.cache.Add(url, dat)
	return locationsResp, nil
}

//Get locations
func (c *Client) GetLocation(areaName string) (CurrentLocation, error) {
	url := baseURL + "/location-area/" + areaName

	if val, ok := c.cache.Get(url); ok {
		locationsResp := CurrentLocation{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			return CurrentLocation{}, err
		}

		return locationsResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CurrentLocation{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CurrentLocation{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return CurrentLocation{}, fmt.Errorf("%d", resp.StatusCode)
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return CurrentLocation{}, err
	}

	locationsResp := CurrentLocation{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return CurrentLocation{}, err
	}

	c.cache.Add(url, dat)
	return locationsResp, nil
}

func (c *Client) GetPokemon(pokemon string) (CurrentPokemon, error) {
	url := baseURL + "/pokemon/" + pokemon

	if val, ok := c.cache.Get(url); ok {
		locationsResp := CurrentPokemon{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			return CurrentPokemon{}, err
		}

		return locationsResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CurrentPokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CurrentPokemon{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode > 299 {
		return CurrentPokemon{}, fmt.Errorf("%d", resp.StatusCode)
	} 

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return CurrentPokemon{}, err
	}

	locationsResp := CurrentPokemon{}
	err = json.Unmarshal(dat, &locationsResp)
	if err != nil {
		return CurrentPokemon{}, err
	}

	c.cache.Add(url, dat)
	return locationsResp, nil
}

func resizeImage(img image.Image, targetWidth, targetHeight int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			// Nearest-neighbor coordinate lookup
			srcX := srcBounds.Min.X + (x * srcW / targetWidth)
			srcY := srcBounds.Min.Y + (y * srcH / targetHeight)
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}

func renderImageLines(img image.Image) []string {
	bounds := img.Bounds()
	height := bounds.Dy()

	var lines []string

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		line := ""
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			topColor := img.At(x, y)
			var bottomColor color.Color = color.RGBA{0, 0, 0, 0}
			if y+1 < height {
				bottomColor = img.At(x, y+1)
			}
			line += formatHalfBlock(topColor, bottomColor)
		}
		line += "\033[0m" // reset formatting at line end
		lines = append(lines, line)
	}

	return lines
}

func formatHalfBlock(top, bottom color.Color) string {
	tr, tg, tb, ta := top.RGBA()
	br, bg, bb, ba := bottom.RGBA()

	topR, topG, topB := uint8(tr>>8), uint8(tg>>8), uint8(tb>>8)
	botR, botG, botB := uint8(br>>8), uint8(bg>>8), uint8(bb>>8)

	topVisible := ta > 0x8000
	botVisible := ba > 0x8000

	switch {
	case topVisible && botVisible:
		return fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀", topR, topG, topB, botR, botG, botB)
	case topVisible && !botVisible:
		return fmt.Sprintf("\033[38;2;%d;%d;%dm\033[49m▀", topR, topG, topB)
	case !topVisible && botVisible:
		return fmt.Sprintf("\033[38;2;%d;%d;%dm\033[49m▄", botR, botG, botB)
	default:
		return "\033[0m "
	}
}

func (c *Client) GetSpriteLines(imageURL string, width, height int) []string {
	if imageURL == "" {
		return nil
	}

	var data []byte
	cached, ok := c.cache.Get(imageURL)
	if ok {
		data = cached
	} else {
		req, err := http.NewRequest("GET", imageURL, nil)
		if err != nil {
			return nil
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 {
			return nil
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		c.cache.Add(imageURL, data)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}

	resized := resizeImage(img, width, height)
	return renderImageLines(resized)
}