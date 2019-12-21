package drawille

type Option func(cfg *Config)

func SaturateOnOverwrite() Option {
	return func(cfg *Config) {
		cfg.overwrite = overwriteSaturate
	}
}

func ToggleOnOverwrite() Option {
	return func(cfg *Config) {
		cfg.overwrite = overwriteToggle
	}
}

func SetPalette(p Palette) Option {
	return func(cfg *Config) {
		cfg.pallette = p
	}
}

type overwriteBehavior int
const (
	overwriteToggle  overwriteBehavior = iota
	overwriteSaturate
)

type Palette map[int]func(string) string

type Config struct {
	overwrite overwriteBehavior
	pallette Palette
}
