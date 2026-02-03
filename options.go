package credui

// Option configures a Client.
type Option func(*clientConfig)

type clientConfig struct {
	mode   Mode
	logger Logger
}

func defaultConfig() clientConfig {
	return clientConfig{
		mode:   ModeProd,
		logger: nopLogger{},
	}
}

// WithMode sets the logging mode.
func WithMode(mode Mode) Option {
	return func(cfg *clientConfig) {
		cfg.mode = mode
	}
}

// WithStdLogger enables the built-in stdout logger.
func WithStdLogger(mode Mode) Option {
	return func(cfg *clientConfig) {
		cfg.mode = mode
		cfg.logger = newStdLogger(mode)
	}
}

// WithLogger sets a custom logger implementation.
func WithLogger(logger Logger) Option {
	return func(cfg *clientConfig) {
		if logger != nil {
			cfg.logger = logger
		}
	}
}
