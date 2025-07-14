package imgtool

import (
	"fmt"
	"imgtool/pkg/converter"
	"imgtool/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().String("to", "png", "format to convert to (png, jpg, gif, webp)")
	convertCmd.Flags().Int("quality", 80, "image quality for jpg/webp (1-100)")
	convertCmd.Flags().String("resize", "", `resize dimensions (e.g., "800x600" or "50%")`)
	convertCmd.Flags().String("watermark", "", "path to watermark image")
	convertCmd.Flags().Int("watermark-opacity", 100, "watermark opacity (0-100)")
	convertCmd.Flags().String("mode", "single", "processing mode (single|dir)")
	convertCmd.Flags().Bool("recursive", false, "recursively process directories")
	convertCmd.Flags().Int("workers", 4, "number of concurrent workers for directory processing")

	viper.BindPFlag("to", convertCmd.Flags().Lookup("to"))
	viper.BindPFlag("quality", convertCmd.Flags().Lookup("quality"))
	viper.BindPFlag("resize", convertCmd.Flags().Lookup("resize"))
	viper.BindPFlag("watermark", convertCmd.Flags().Lookup("watermark"))
	viper.BindPFlag("watermark-opacity", convertCmd.Flags().Lookup("watermark-opacity"))
	viper.BindPFlag("workers", convertCmd.Flags().Lookup("workers"))
}

var convertCmd = &cobra.Command{
	Use:   "convert [input_path] [output_path]",
	Short: "Convert and process images",
	Long:  `Convert images to different formats, resize, adjust quality, and apply watermarks. Supports single file and directory batch processing.`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]
		outputPath := ""
		if len(args) > 1 {
			outputPath = args[1]
		}

		mode, _ := cmd.Flags().GetString("mode")
		recursive, _ := cmd.Flags().GetBool("recursive")

		opts := converter.Options{
			To:               viper.GetString("to"),
			Quality:          viper.GetInt("quality"),
			Resize:           viper.GetString("resize"),
			WatermarkPath:    viper.GetString("watermark"),
			WatermarkOpacity: viper.GetInt("watermark-opacity"),
		}

		if mode == "single" {
			if outputPath == "" {
				outputPath = generateOutputFilename(inputPath, opts.To)
			}
			fmt.Printf("Processing %s -> %s\n", inputPath, outputPath)
			err := converter.ProcessImage(inputPath, outputPath, opts)
			if err != nil {
				log.Fatalf("Failed to process image: %v", err)
			}
		} else if mode == "dir" {
			if outputPath == "" {
				log.Fatal("Output directory is required for 'dir' mode")
			}
			os.MkdirAll(outputPath, os.ModePerm)

			workers := viper.GetInt("workers")
			files, err := utils.FindImageFiles(inputPath, recursive)
			if err != nil {
				log.Fatalf("Failed to find image files: %v", err)
			}

			converter.ProcessImages(files, outputPath, opts, workers)
		} else {
			log.Fatalf("Invalid mode: %s. Use 'single' or 'dir'", mode)
		}

		fmt.Println("Image processing complete.")
	},
}

func generateOutputFilename(inputPath, toFormat string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return filepath.Join(dir, fmt.Sprintf("%s.%s", name, toFormat))
}
