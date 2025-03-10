package gci

import (
	"fmt"

	gcicfg "github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/log"
	"github.com/ldez/grignotin/gomod"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goformatters/internal"
)

const Name = "gci"

type Formatter struct {
	config *gcicfg.Config
}

func New(settings *config.GciSettings) (*Formatter, error) {
	log.InitLogger()
	_ = log.L().Sync()

	modPath, err := gomod.GetModulePath()
	if err != nil {
		internal.FormatterLogger.Errorf("gci: %v", err)
	}

	cfg := gcicfg.YamlConfig{
		Cfg: gcicfg.BoolConfig{
			NoInlineComments: settings.NoInlineComments,
			NoPrefixComments: settings.NoPrefixComments,
			SkipGenerated:    settings.SkipGenerated,
			CustomOrder:      settings.CustomOrder,
			NoLexOrder:       settings.NoLexOrder,
		},
		SectionStrings: settings.Sections,
		ModPath:        modPath,
	}

	if settings.LocalPrefixes != "" {
		cfg.SectionStrings = []string{
			"standard",
			"default",
			fmt.Sprintf("prefix(%s)", settings.LocalPrefixes),
		}
	}

	parsedCfg, err := cfg.Parse()
	if err != nil {
		return nil, err
	}

	return &Formatter{config: parsedCfg}, nil
}

func (*Formatter) Name() string {
	return Name
}

func (f *Formatter) Format(filename string, src []byte) ([]byte, error) {
	_, formatted, err := gci.LoadFormat(src, filename, *f.config)
	return formatted, err
}
