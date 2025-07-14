# imgtool

A powerful and fast CLI tool for image processing, built with Go, Cobra, and Viper.

**Link**: [https://github.com/PrathamGhaywat/imgtool](https://github.com/PrathamGhaywat/imgtool)

---

## Features

- **Format Conversion**: Convert images between `png`, `jpg`, `gif`, and `webp`.
- **Quality Adjustment**: Set the quality for lossy formats like `jpg` and `webp`.
- **Image Resizing**: Resize images using specific dimensions (e.g., `800x600`) or by a percentage (e.g., `50%`). Aspect ratio is preserved if one dimension is set to `0`.
- **Watermarking**: Apply a watermark image with adjustable opacity.
- **Batch Processing**: Process a single file or an entire directory of images.
- **Concurrency**: Use multiple workers for faster directory processing.
- **Configuration**: Set default flag values using a `config.yaml` file.

## Installation
### Windows
Install the release from [here](https://github.com/PrathamGhaywat/imgtool/releases). Then add the folder in which the executable is present in your `PATH` environment variable.

### Linux and macOS
We are currently searching for a way to make the tool available for Linux and macOS. Please open an issue if you have any suggestions.
If you would still like to use the tool on Linux and macOS, you can [build](#build) it from the source.

## Documentation
### Documentation

The documentation for imgtool can be found [here](https://prathamghaywat.github.io/imgtool-docs/).


## Build
To build the tool from the source, you'll need to have Go installed.

```bash
# Clone the repository
git clone https://github.com/PrathamGhaywat/imgtool.git
cd imgtool

# Download dependencies
go mod tidy

# Build the executable
go build -o imgtool.exe .
```

## Usage

### General Help

To see a list of all available commands and flags:

```bash
.\imgtool.exe --help
```

### Convert Command

The `convert` command is the primary tool for all image processing tasks.

```bash
.\imgtool.exe convert [input_path] [output_path] [flags]
```

#### Arguments

- `input_path`: The path to the source image or directory.
- `output_path` (optional): The path for the output file or directory. If omitted for a single file, the output will be saved in the same directory with the new format extension.

#### Flags

- `--to <format>`: The target format (`png`, `jpg`, `gif`, `webp`). Default: `png`.
- `--quality <int>`: Image quality for `jpg`/`webp` (1-100). Default: `80`.
- `--resize <string>`: Resize dimensions (e.g., `"800x600"` or `"50%"`).
- `--watermark <path>`: Path to the watermark image.
- `--watermark-opacity <int>`: Watermark opacity (0-100). Default: `100`.
- `--mode <string>`: Processing mode (`single` or `dir`). Default: `single`.
- `--recursive`: Recursively process subdirectories in `dir` mode.
- `--workers <int>`: Number of concurrent workers for `dir` mode. Default: `4`.

### Examples

1.  **Convert a PNG to a high-quality JPG**

    ```bash
    .\imgtool.exe convert input.png output.jpg --to jpg --quality 95
    ```

2.  **Resize an image to 50% of its original size**

    ```bash
    .\imgtool.exe convert big.jpg small.jpg --resize "50%"
    ```

3.  **Apply a watermark to an image**

    ```bash
    .\imgtool.exe convert photo.png watermarked.png --watermark logo.png --watermark-opacity 70
    ```

4.  **Process an entire directory of images**

    This command converts all images in the `input-folder` to `webp` and saves them in `output-folder`.

    ```bash
    .\imgtool.exe convert input-folder output-folder --mode dir --to webp
    ```

## Configuration

You can set default values for any flag in the `config.yaml` file. The tool looks for this file in the current directory or your home directory.

**Example `config.yaml`:**

```yaml
# Default configuration for imgtool
to: "webp"
quality: 75
watermark: "C:/Users/YourUser/Pictures/watermark.png"
watermark-opacity: 60
workers: 8
```
