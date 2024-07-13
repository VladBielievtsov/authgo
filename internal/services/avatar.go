package services

import (
	"authgo/internal/utils"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var colors = []color.RGBA{
	{25, 165, 255, 255},
	{254, 79, 25, 255},
	{240, 12, 24, 255},
	{255, 23, 145, 255},
	{4, 231, 98, 255},
	{137, 0, 242, 255},
}

type AvatarServices struct{}

func NewAvatarServices() *AvatarServices {
	return &AvatarServices{}
}

func (s *AvatarServices) GenerateAvatar(firstName, id string) (string, error) {
	width := 200
	height := 200
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	fillColor := colors[rand.Intn(len(colors))]
	draw.Draw(img, img.Bounds(), &image.Uniform{fillColor}, image.Point{}, draw.Src)

	fontPath := filepath.Join("fonts", "Roboto-Regular.ttf")

	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return "", fmt.Errorf("failed to read font file: %v", err)
	}

	font, err := opentype.Parse(fontBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse font: %v", err)
	}

	fontSize := 120.0
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: fontSize,
		DPI:  72,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create new face: %v", err)
	}
	defer face.Close()

	initial := strings.ToUpper(string(firstName[0]))
	err = addLabel(img, width/2, height/2, initial, face)
	if err != nil {
		return "", fmt.Errorf("failed to add label: %v", err)
	}

	if err := utils.EnsureDirExists(filepath.Join("uploads", "avatars")); err != nil {
		return "", fmt.Errorf("failed to create uploads/avatars directory: %v", err)
	}

	avatarPath := filepath.Join("uploads", "avatars")
	filename := fmt.Sprintf("%s/%s.png", avatarPath, id)

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %v", err)
	}

	return filename, nil
}

func addLabel(img *image.RGBA, x, y int, label string, face font.Face) error {
	col := color.RGBA{255, 255, 255, 255}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}

	bounds, _ := d.BoundString(label)
	d.Dot.X -= (bounds.Max.X - bounds.Min.X) / 2
	d.Dot.Y += (bounds.Max.Y - bounds.Min.Y) / 2

	d.DrawString(label)

	return nil
}
