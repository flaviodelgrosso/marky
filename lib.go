package marky

import (
	"github.com/flaviodelgrosso/marky/internal/converters"
	"github.com/flaviodelgrosso/marky/internal/marky"
)

// Creates a new marky instance with all available loaders registered.
func New() marky.IMarky {
	m := &marky.Marky{
		Converters: make([]converters.Converter, 0, 8),
	}

	m.RegisterConverter(converters.NewCsvConverter())
	m.RegisterConverter(converters.NewDocConverter())
	m.RegisterConverter(converters.NewEpubConverter())
	m.RegisterConverter(converters.NewExcelConverter())
	m.RegisterConverter(converters.NewHTMLConverter())
	m.RegisterConverter(converters.NewIpynbConverter())
	m.RegisterConverter(converters.NewPdfConverter())
	m.RegisterConverter(converters.NewPptxConverter())

	return m
}
