package credui

// Client provides methods to show credential dialogs.
//
// On non-Windows platforms, all methods return ErrUnsupportedPlatform.
type Client struct {
	cfg clientConfig
	api platformAPI
}

// New creates a new Client.
func New(opts ...Option) *Client {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	cfg.logger.Debugf("credui.New: mode=%v logger=%T", cfg.mode, cfg.logger)

	return &Client{
		cfg: cfg,
		api: newPlatformAPI(cfg.logger),
	}
}

// PromptWindowsCredentials shows the modern credential dialog.
func (c *Client) PromptWindowsCredentials(opts WindowsPromptOptions) (Credential, error) {
	c.cfg.logger.Debugf("Client.PromptWindowsCredentials: opts=%+v", opts)
	cred, err := c.api.PromptWindowsCredentials(opts)
	c.cfg.logger.Debugf("Client.PromptWindowsCredentials: result=%+v err=%v", redacted(cred), err)
	return cred, err
}

// PromptClassic shows the legacy credential dialog.
func (c *Client) PromptClassic(opts ClassicPromptOptions) (Credential, error) {
	c.cfg.logger.Debugf("Client.PromptClassic: opts=%+v", opts)
	cred, err := c.api.PromptClassic(opts)
	c.cfg.logger.Debugf("Client.PromptClassic: result=%+v err=%v", redacted(cred), err)
	return cred, err
}

func redacted(c Credential) Credential {
	out := c
	if out.Password != "" {
		out.Password = "***"
	}
	return out
}
