# Marky 📝

[![Go Version](https://img.shields.io/github/go-mod/go-version/flaviodelgrosso/marky)](https://golang.org/doc/go1.24)
[![License](https://img.shields.io/badge/license-ISC-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/flaviodelgrosso/marky)](https://goreportcard.com/report/github.com/flaviodelgrosso/marky)

A powerful Go library and CLI tool for converting various document formats to Markdown. Marky makes it easy to extract and convert content from different file types into clean, readable Markdown format.

## 🚀 Features

- **Multiple Format Support**: Convert CSV, HTML, Jupiter Notebooks, Word, Excel, PDF, and PowerPoint files to Markdown
- **CLI Tool**: Easy-to-use command-line interface for quick conversions
- **Go Library**: Integrate conversion capabilities into your Go applications
- **MIME Type Detection**: Automatic file type detection for robust handling
- **Extensible Architecture**: Plugin-based loader system for easy format additions

## 📋 Supported Formats

| Format | Extensions | MIME Types |
|--------|------------|------------|
| **CSV** | `.csv` | `text/csv`, `application/csv` |
| **HTML** | `.html`, `.htm` | `text/html` |
| **Jupyter Notebook** | `.ipynb` | `application/x-ipynb+json`, `application/json` |
| **Microsoft Word** | `.docx` | `application/vnd.openxmlformats-officedocument.wordprocessingml.document` |
| **Microsoft Excel** | `.xlsx` | `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` |
| **PDF** | `.pdf` | `application/pdf` |
| **Microsoft PowerPoint** | `.pptx` | `application/vnd.openxmlformats-officedocument.presentationml.presentation` |

## 📦 Installation

### CLI Tool

Install the CLI tool directly using Go:

```bash
go install github.com/flaviodelgrosso/marky/cmd/marky@latest
```

### Library

Add Marky to your Go project:

```bash
go get github.com/flaviodelgrosso/marky
```

## 🛠️ Usage

### Command Line Interface

Basic usage:

```bash
# Convert a file and output to console
marky document.pdf

# Convert a file and save to output file
marky document.docx --output converted.md
marky document.xlsx -o converted.md

# Examples with different formats
marky presentation.pptx -o slides.md
marky data.csv -o table.md
marky webpage.html -o content.md
```

### Go Library

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/flaviodelgrosso/marky"
)

func main() {
    // Initialize Marky with all available loaders
    m := marky.New()
    
    // Convert a document to Markdown
    result, err := m.Convert("document.pdf")
    if err != nil {
        log.Fatalf("Conversion failed: %v", err)
    }
    
    fmt.Println(result)
}
```

## 🏗️ Development

### Prerequisites

- Go 1.24.4 or later
- Make (optional, for using Makefile commands)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/flaviodelgrosso/marky.git
cd marky

# Build the CLI tool
make build
# OR
go build -o bin/marky cmd/marky/main.go

# Run tests
make test
# OR
go test -v ./...

# Run linting (requires golangci-lint)
make lint
```

## 🧪 Testing

Run the test suite:

```bash
go test -v ./...
```

Test files for various formats are included in the `test_files/` directory to ensure proper functionality across all supported document types.

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** and add tests
4. **Run tests**: `make test`
5. **Commit your changes**: `git commit -m 'Add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Adding New Format Support

To add support for a new document format:

1. Create a new loader in `internal/loaders/`
2. Implement the `DocumentLoader` interface:

   ```go
   type DocumentLoader interface {
       Load(path string) (string, error)
       CanLoadMimeType(mimeType string) bool
   }
   ```

3. Register the loader in the `New()` function in `lib.go`
4. Add tests for your new loader

## 📄 License

This project is licensed under the ISC License. See the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) for HTML conversion
- [pdf](https://github.com/ledongthuc/pdf) for PDF text extraction
- [excelize](https://github.com/xuri/excelize) for Excel file processing
- [cobra](https://github.com/spf13/cobra) for CLI framework

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/flaviodelgrosso/marky/issues)
- 💡 **Feature Requests**: [GitHub Issues](https://github.com/flaviodelgrosso/marky/issues)
- 📧 **Questions**: Open a [GitHub Discussion](https://github.com/flaviodelgrosso/marky/discussions)

---

Made with ❤️ by [Flavio Del Grosso](https://github.com/flaviodelgrosso)
