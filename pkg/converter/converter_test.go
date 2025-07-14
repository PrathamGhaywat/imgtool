package converter

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createDummyImage creates a simple image for testing.
func createDummyImage(t *testing.T, width, height int, name string) string {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{100, 200, 200, 255})
		}
	}

	tmpFile, err := os.CreateTemp("", name+"*.png")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	if err := png.Encode(tmpFile, img); err != nil {
		t.Fatalf("Failed to encode dummy image: %v", err)
	}

	return tmpFile.Name()
}

func TestProcessImage_Convert(t *testing.T) {
	inputPath := createDummyImage(t, 10, 10, "input")
	defer os.Remove(inputPath)

	outputPath := filepath.Join(t.TempDir(), "output.jpg")
	opts := Options{
		To:      "jpg",
		Quality: 90,
	}

	err := ProcessImage(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("ProcessImage() failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created")
	}
}

func TestProcessImage_Resize(t *testing.T) {
	inputPath := createDummyImage(t, 200, 200, "resize-input")
	defer os.Remove(inputPath)

	outputPath := filepath.Join(t.TempDir(), "output-resized.png")
	opts := Options{
		To:     "png",
		Resize: "100x100",
	}

	err := ProcessImage(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("ProcessImage() with resize failed: %v", err)
	}

	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("Failed to open resized image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("Failed to decode resized image: %v", err)
	}

	if img.Bounds().Dx() != 100 || img.Bounds().Dy() != 100 {
		t.Errorf("Expected image dimensions 100x100, got %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestProcessImage_Watermark(t *testing.T) {
	basePath := createDummyImage(t, 300, 300, "base")
	defer os.Remove(basePath)

	watermarkPath := createDummyImage(t, 50, 50, "watermark")
	defer os.Remove(watermarkPath)

	outputPath := filepath.Join(t.TempDir(), "output-watermarked.png")
	opts := Options{
		To:               "png",
		WatermarkPath:    watermarkPath,
		WatermarkOpacity: 80,
	}

	err := ProcessImage(basePath, outputPath, opts)
	if err != nil {
		t.Fatalf("ProcessImage() with watermark failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Watermarked output file was not created")
	}
}
