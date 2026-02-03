package credui

// platformAPI abstracts OS-specific implementations.
type platformAPI interface {
	PromptWindowsCredentials(opts WindowsPromptOptions) (Credential, error)
	PromptClassic(opts ClassicPromptOptions) (Credential, error)
}
