package credui

import "errors"

var (
	// ErrUnsupportedPlatform indicates the API is not available on this OS.
	ErrUnsupportedPlatform = errors.New("credui: unsupported platform")
	// ErrCancelled indicates the user cancelled the dialog.
	ErrCancelled = errors.New("credui: cancelled by user")
)
