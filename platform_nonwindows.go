//go:build !windows

package credui

type unsupportedAPI struct{ logger Logger }

func newPlatformAPI(logger Logger) platformAPI {
	if logger != nil {
		logger.Debugf("newPlatformAPI: non-windows -> unsupportedAPI")
	}
	return &unsupportedAPI{logger: logger}
}

func (u *unsupportedAPI) PromptWindowsCredentials(opts WindowsPromptOptions) (Credential, error) {
	if u.logger != nil {
		u.logger.Debugf("unsupportedAPI.PromptWindowsCredentials: opts=%+v -> ErrUnsupportedPlatform", opts)
	}
	return Credential{}, ErrUnsupportedPlatform
}

func (u *unsupportedAPI) PromptClassic(opts ClassicPromptOptions) (Credential, error) {
	if u.logger != nil {
		u.logger.Debugf("unsupportedAPI.PromptClassic: opts=%+v -> ErrUnsupportedPlatform", opts)
	}
	return Credential{}, ErrUnsupportedPlatform
}
