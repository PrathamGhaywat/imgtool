package converter

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/chai2010/webp"
	xdraw "golang.org/x/image/draw"
)

// Options holds the processing options for an image.
type Options struct {
	To               string
	Quality          int
	Resize           string
	WatermarkPath    string
	WatermarkOpacity int
}

// ProcessImage handles the full conversion pipeline for a single image.
func ProcessImage(inputPath, outputPath string, opts Options) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	if opts.Resize != "" {
		img = resize(img, opts.Resize)
	}

	if opts.WatermarkPath != "" {
		img, err = applyWatermark(img, opts.WatermarkPath, opts.WatermarkOpacity)
		if err != nil {
			return fmt.Errorf("failed to apply watermark: %w", err)
		}
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	switch strings.ToLower(opts.To) {
	case "jpg", "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: opts.Quality})
	case "png":
		return png.Encode(out, img)
	case "gif":
		return gif.Encode(out, img, &gif.Options{NumColors: 256})
	case "webp":
		return webp.Encode(out, img, &webp.Options{Lossless: false, Quality: float32(opts.Quality)})
	default:
		return fmt.Errorf("unsupported output format: %s", opts.To)
	}
}

// ProcessImages processes a list of images concurrently.
func ProcessImages(files []string, outputDir string, opts Options, workers int) {
	jobs := make(chan string, len(files))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for inputPath := range jobs {
				base := filepath.Base(inputPath)
				ext := filepath.Ext(base)
				name := strings.TrimSuffix(base, ext)
				outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.%s", name, opts.To))

				fmt.Printf("Processing %s -> %s\n", inputPath, outputPath)
				if err := ProcessImage(inputPath, outputPath, opts); err != nil {
					log.Printf("Failed to process %s: %v", inputPath, err)
				}
			}
		}()
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	wg.Wait()
}

func resize(img image.Image, resizeStr string) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var newWidth, newHeight int

	if strings.HasSuffix(resizeStr, "%") {
		pctStr := strings.TrimSuffix(resizeStr, "%")
		pct, err := strconv.Atoi(pctStr)
		if err != nil {
			log.Printf("Invalid resize percentage: %s", resizeStr)
			return img
		}
		newWidth = width * pct / 100
		newHeight = height * pct / 100
	} else {
		parts := strings.Split(strings.ToLower(resizeStr), "x")
		if len(parts) != 2 {
			log.Printf("Invalid resize format: %s", resizeStr)
			return img
		}
		newWidth, _ = strconv.Atoi(parts[0])
		newHeight, _ = strconv.Atoi(parts[1])
	}

	if newWidth == 0 || newHeight == 0 {
		// Preserve aspect ratio
		aspectRatio := float64(width) / float64(height)
		if newWidth == 0 {
			newWidth = int(float64(newHeight) * aspectRatio)
		} else {
			newHeight = int(float64(newWidth) / aspectRatio)
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	return dst
}

func applyWatermark(baseImg image.Image, watermarkPath string, opacity int) (image.Image, error) {
	wmFile, err := os.Open(watermarkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open watermark image: %w", err)
	}
	defer wmFile.Close()

	wmImg, _, err := image.Decode(wmFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode watermark image: %w", err)
	}

	offset := image.Pt(baseImg.Bounds().Dx()-wmImg.Bounds().Dx()-10, baseImg.Bounds().Dy()-wmImg.Bounds().Dy()-10)
	b := baseImg.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, b, baseImg, image.Point{}, draw.Src)

	mask := image.NewUniform(color.Alpha{A: uint8(float64(opacity) / 100.0 * 255.0)})
	draw.DrawMask(dst, wmImg.Bounds().Add(offset), wmImg, image.Point{}, mask, image.Point{}, draw.Over)

	return dst, nil
}
