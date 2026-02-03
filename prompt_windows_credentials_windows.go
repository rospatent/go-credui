//go:build windows

package credui

import (
	"fmt"
	"syscall"
	"unsafe"
)

// PromptWindowsCredentials shows the modern credential prompt using
// CredUIPromptForWindowsCredentialsW.
func (w *winAPIs) PromptWindowsCredentials(opts WindowsPromptOptions) (Credential, error) {
	if w.logger != nil {
		w.logger.Debugf("winAPIs.PromptWindowsCredentials: start opts=%+v", opts)
	}

	caption, err := utf16PtrOrNil(opts.Caption)
	if err != nil {
		if w.logger != nil {
			w.logger.Errorf("winAPIs.PromptWindowsCredentials: bad Caption=%q err=%v", opts.Caption, err)
		}
		return Credential{}, fmt.Errorf("caption UTF16 conversion: %w", err)
	}
	message, err := utf16PtrOrNil(opts.Message)
	if err != nil {
		if w.logger != nil {
			w.logger.Errorf("winAPIs.PromptWindowsCredentials: bad Message=%q err=%v", opts.Message, err)
		}
		return Credential{}, fmt.Errorf("message UTF16 conversion: %w", err)
	}

	ui := credUIInfo{
		CbSize:     uint32(unsafe.Sizeof(credUIInfo{})),
		HwndParent: opts.ParentHWND,
		PszCaption: caption,
		PszMessage: message,
	}

	flags := opts.Flags
	if flags == 0 {
		flags = CREDUIWIN_GENERIC
	}

	var authPkg uint32
	var outBuf uintptr
	var outSize uint32
	var saveCred bool

	if w.logger != nil {
		w.logger.Debugf("winAPIs.PromptWindowsCredentials: calling CredUIPromptForWindowsCredentialsW flags=0x%x", flags)
	}

	r, _, _ := w.promptWindowsCreds.Call(
		uintptr(unsafe.Pointer(&ui)),
		0, // dwAuthError
		uintptr(unsafe.Pointer(&authPkg)),
		0, // pvInAuthBuffer
		0, // ulInAuthBufferSize
		uintptr(unsafe.Pointer(&outBuf)),
		uintptr(unsafe.Pointer(&outSize)),
		uintptr(unsafe.Pointer(&saveCred)),
		uintptr(flags),
	)
	if err := winCallError("CredUIPromptForWindowsCredentialsW", r); err != nil {
		if w.logger != nil {
			w.logger.Warnf("winAPIs.PromptWindowsCredentials: prompt returned err=%v", err)
		}
		return Credential{}, err
	}
	if outBuf == 0 || outSize == 0 {
		if w.logger != nil {
			w.logger.Errorf("winAPIs.PromptWindowsCredentials: outBuf/outSize empty outBuf=0x%x outSize=%d", outBuf, outSize)
		}
		return Credential{}, fmt.Errorf("CredUIPromptForWindowsCredentialsW returned empty buffer")
	}
	defer w.freeCoTaskMem(outBuf)

	if w.logger != nil {
		w.logger.Debugf("winAPIs.PromptWindowsCredentials: unpack pass #1 (sizes) outBuf=0x%x outSize=%d", outBuf, outSize)
	}

	var userLen, domainLen, passLen uint32
	// First call: ask for required sizes.
	r2, _, _ := w.unpackAuthBuffer.Call(
		0, // dwFlags
		outBuf,
		uintptr(outSize),
		0,
		uintptr(unsafe.Pointer(&userLen)),
		0,
		uintptr(unsafe.Pointer(&domainLen)),
		0,
		uintptr(unsafe.Pointer(&passLen)),
	)
	// Expect failure with ERROR_INSUFFICIENT_BUFFER in many cases; the sizes are still useful.
	if w.logger != nil {
		w.logger.Debugf("winAPIs.PromptWindowsCredentials: unpack #1 r=%d userLen=%d domainLen=%d passLen=%d", r2, userLen, domainLen, passLen)
	}

	if userLen == 0 {
		// Some providers may return empty userLen on first call; allocate a reasonable default.
		userLen = 256
	}
	if domainLen == 0 {
		domainLen = 256
	}
	if passLen == 0 {
		passLen = 256
	}

	userBuf := make([]uint16, userLen)
	domainBuf := make([]uint16, domainLen)
	passBuf := make([]uint16, passLen)
	defer zeroUint16Slice(passBuf)

	if w.logger != nil {
		w.logger.Debugf("winAPIs.PromptWindowsCredentials: unpack pass #2 buffers user=%d domain=%d pass=%d", len(userBuf), len(domainBuf), len(passBuf))
	}

	r3, _, _ := w.unpackAuthBuffer.Call(
		0,
		outBuf,
		uintptr(outSize),
		uintptr(unsafe.Pointer(&userBuf[0])),
		uintptr(unsafe.Pointer(&userLen)),
		uintptr(unsafe.Pointer(&domainBuf[0])),
		uintptr(unsafe.Pointer(&domainLen)),
		uintptr(unsafe.Pointer(&passBuf[0])),
		uintptr(unsafe.Pointer(&passLen)),
	)
	if r3 == 0 {
		// CredUnPackAuthenticationBufferW returns BOOL (0==fail).
		lastErr := syscall.GetLastError()
		if w.logger != nil {
			w.logger.Errorf("winAPIs.PromptWindowsCredentials: CredUnPackAuthenticationBufferW failed lastErr=%v", lastErr)
		}
		return Credential{}, fmt.Errorf("CredUnPackAuthenticationBufferW failed: %w", lastErr)
	}

	cred := Credential{
		Username:     utf16SliceToString(userBuf),
		Domain:       utf16SliceToString(domainBuf),
		Password:     utf16SliceToString(passBuf),
		SaveSelected: saveCred,
	}

	if w.logger != nil {
		w.logger.Debugf("winAPIs.PromptWindowsCredentials: success username=%q domain=%q save=%v passLen=%d", cred.Username, cred.Domain, cred.SaveSelected, len(cred.Password))
	}
	return cred, nil
}
