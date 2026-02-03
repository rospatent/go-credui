//go:build windows

package credui

import (
	"fmt"
	"unsafe"
)

// PromptClassic shows the classic credential prompt using CredUIPromptForCredentialsW.
func (w *winAPIs) PromptClassic(opts ClassicPromptOptions) (Credential, error) {
	if w.logger != nil {
		w.logger.Debugf("winAPIs.PromptClassic: start opts=%+v", opts)
	}

	targetName := opts.TargetName
	if targetName == "" {
		targetName = "CREDENTIAL"
	}
	targetPtr, err := utf16PtrOrNil(targetName)
	if err != nil {
		if w.logger != nil {
			w.logger.Errorf("winAPIs.PromptClassic: bad TargetName=%q err=%v", targetName, err)
		}
		return Credential{}, fmt.Errorf("target UTF16 conversion: %w", err)
	}

	caption, err := utf16PtrOrNil(opts.Caption)
	if err != nil {
		return Credential{}, fmt.Errorf("caption UTF16 conversion: %w", err)
	}
	message, err := utf16PtrOrNil(opts.Message)
	if err != nil {
		return Credential{}, fmt.Errorf("message UTF16 conversion: %w", err)
	}

	var uiPtr uintptr
	var ui credUIInfo
	if caption != nil || message != nil || opts.ParentHWND != 0 {
		ui = credUIInfo{
			CbSize:     uint32(unsafe.Sizeof(credUIInfo{})),
			HwndParent: opts.ParentHWND,
			PszCaption: caption,
			PszMessage: message,
		}
		uiPtr = uintptr(unsafe.Pointer(&ui))
	}

	flags := opts.Flags
	if flags == 0 {
		flags = CREDUI_FLAGS_GENERIC_CREDENTIALS | CREDUI_FLAGS_DO_NOT_PERSIST
	}

	userMax := opts.UsernameMaxChars
	if userMax == 0 {
		userMax = 256
	}
	passMax := opts.PasswordMaxChars
	if passMax == 0 {
		passMax = 256
	}

	userBuf := make([]uint16, userMax)
	passBuf := make([]uint16, passMax)
	defer zeroUint16Slice(passBuf)

	var save bool

	if w.logger != nil {
		w.logger.Debugf("winAPIs.PromptClassic: calling CredUIPromptForCredentialsW target=%q flags=0x%x userMax=%d passMax=%d", targetName, flags, userMax, passMax)
	}

	r, _, _ := w.promptClassic.Call(
		uiPtr,
		uintptr(unsafe.Pointer(targetPtr)),
		0, // pvReserved (PCtxtHandle context)
		0, // dwAuthError
		uintptr(unsafe.Pointer(&userBuf[0])),
		uintptr(userMax),
		uintptr(unsafe.Pointer(&passBuf[0])),
		uintptr(passMax),
		uintptr(unsafe.Pointer(&save)),
		uintptr(flags),
	)
	if err := winCallError("CredUIPromptForCredentialsW", r); err != nil {
		if w.logger != nil {
			w.logger.Warnf("winAPIs.PromptClassic: prompt returned err=%v", err)
		}
		return Credential{}, err
	}

	cred := Credential{
		Username:     utf16SliceToString(userBuf),
		Password:     utf16SliceToString(passBuf),
		SaveSelected: save,
	}
	if w.logger != nil {
		w.logger.Debugf("winAPIs.PromptClassic: success username=%q save=%v passLen=%d", cred.Username, cred.SaveSelected, len(cred.Password))
	}
	return cred, nil
}
