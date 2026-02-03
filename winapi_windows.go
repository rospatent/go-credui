//go:build windows

package credui

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	errorCancelled = 1223 // ERROR_CANCELLED

	// CredUIPromptForWindowsCredentials dwFlags.
	CREDUIWIN_GENERIC  = 0x00000001
	CREDUIWIN_CHECKBOX = 0x00000002

	// CredUIPromptForCredentials dwFlags.
	CREDUI_FLAGS_GENERIC_CREDENTIALS = 0x00040000
	CREDUI_FLAGS_DO_NOT_PERSIST      = 0x00000002
)

// credUIInfo matches CREDUI_INFO in Win32.
type credUIInfo struct {
	CbSize     uint32
	HwndParent uintptr
	PszMessage *uint16
	PszCaption *uint16
	HbmBanner  uintptr
}

type winAPIs struct {
	logger Logger

	credui *syscall.LazyDLL
	ole32  *syscall.LazyDLL

	promptWindowsCreds *syscall.LazyProc
	unpackAuthBuffer   *syscall.LazyProc
	promptClassic      *syscall.LazyProc
	coTaskMemFree      *syscall.LazyProc
}

func newWinAPIs(logger Logger) *winAPIs {
	api := &winAPIs{logger: logger}
	api.credui = syscall.NewLazyDLL("credui.dll")
	api.ole32 = syscall.NewLazyDLL("ole32.dll")

	api.promptWindowsCreds = api.credui.NewProc("CredUIPromptForWindowsCredentialsW")
	api.unpackAuthBuffer = api.credui.NewProc("CredUnPackAuthenticationBufferW")
	api.promptClassic = api.credui.NewProc("CredUIPromptForCredentialsW")
	api.coTaskMemFree = api.ole32.NewProc("CoTaskMemFree")

	if logger != nil {
		logger.Debugf("newWinAPIs: dlls initialized")
	}
	return api
}

func (w *winAPIs) freeCoTaskMem(ptr uintptr) {
	if ptr == 0 {
		return
	}
	if w.logger != nil {
		w.logger.Debugf("winAPIs.freeCoTaskMem: ptr=0x%x", ptr)
	}
	_, _, _ = w.coTaskMemFree.Call(ptr)
}

func winCallError(fn string, r uintptr) error {
	if r == 0 {
		return nil
	}
	if r == errorCancelled {
		return ErrCancelled
	}
	errno := syscall.Errno(r)
	return fmt.Errorf("%s failed: %w", fn, errno)
}

func utf16PtrOrNil(s string) (*uint16, error) {
	if s == "" {
		return nil, nil
	}
	p, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func utf16SliceToString(s []uint16) string {
	// syscall.UTF16ToString stops on first NUL.
	return syscall.UTF16ToString(s)
}

func zeroUint16Slice(s []uint16) {
	for i := range s {
		s[i] = 0
	}
	// Keep compiler from optimizing the loop away.
	_ = unsafe.Pointer(&s)
}
