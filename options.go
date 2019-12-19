package drawille

type Option func(cfg *Config)

func SaturateOnOverwrite() Option {
	return func(cfg *Config) {
		cfg.owb = overwriteSaturate
	}
}

func ToggleOnOverwrite() Option {
	return func(cfg *Config) {
		cfg.owb = overwriteToggle
	}
}

type overwriteBehavior int
const (
	overwriteToggle  overwriteBehavior = iota
	overwriteSaturate
)

type Config struct {
	owb overwriteBehavior
}
