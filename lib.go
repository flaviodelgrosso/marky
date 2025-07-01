package marky

import (
	"github.com/flaviodelgrosso/marky/internal/loaders"
	"github.com/flaviodelgrosso/marky/internal/marky"
)

// Initialize creates a new marky instance with all available loaders registered.
func Initialize() marky.IMarky {
	m := &marky.Marky{
		Loaders: make([]loaders.DocumentLoader, 0, 7), // Pre-allocate with known capacity
	}

	m.RegisterLoader(&loaders.CsvLoader{})
	m.RegisterLoader(&loaders.DocLoader{})
	m.RegisterLoader(&loaders.ExcelLoader{})
	m.RegisterLoader(&loaders.HTMLLoader{})
	m.RegisterLoader(&loaders.IpynbLoader{})
	m.RegisterLoader(&loaders.PdfLoader{})
	m.RegisterLoader(&loaders.PptxLoader{})

	return m
}
