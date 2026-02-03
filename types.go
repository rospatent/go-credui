package credui

// Credential contains the user-entered credential data.
//
// Security note: returning Password as string means it will live in the Go heap
// until GC. If you need stronger in-memory protection, change the API to return
// a byte slice and explicitly zero it after use.
type Credential struct {
	Username     string
	Domain       string
	Password     string
	SaveSelected bool
}

// WindowsPromptOptions configures the modern dialog (CredUIPromptForWindowsCredentialsW).
type WindowsPromptOptions struct {
	Caption    string
	Message    string
	ParentHWND uintptr
	// Flags are CredUIWIN_* flags.
	// If zero, the library uses CREDUIWIN_GENERIC.
	Flags uint32
}

// ClassicPromptOptions configures the classic dialog (CredUIPromptForCredentialsW).
type ClassicPromptOptions struct {
	TargetName string
	Caption    string
	Message    string
	ParentHWND uintptr
	// Flags are CREDUI_FLAGS_* flags.
	// If zero, the library uses CREDUI_FLAGS_GENERIC_CREDENTIALS | CREDUI_FLAGS_DO_NOT_PERSIST.
	Flags uint32
	// UsernameMaxChars and PasswordMaxChars control the buffer sizes passed to the WinAPI.
	// If zero, defaults are used.
	UsernameMaxChars uint32
	PasswordMaxChars uint32
}
