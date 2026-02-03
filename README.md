# credui (Windows)

A small Go library that wraps Windows Credential UI prompts.

Features:
- Modern UI via `CredUIPromptForWindowsCredentialsW` + `CredUnPackAuthenticationBufferW`.
- Classic UI via `CredUIPromptForCredentialsW`.
- Verbose debug logging if credui.New created with logger.

## Quick start

```go
package main

import (
	"fmt"

	"github.com/rospatent/go-credui"
)

func main() {
	client := credui.New()

	cred, err := client.PromptWindowsCredentials(credui.WindowsPromptOptions{
		Caption: "hi",
		Message: "enter login/password:",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("user:", cred.Username)
	fmt.Println("pass:", cred.Password)
}
```

## Build

This library is Windows-only (the prompts rely on `credui.dll`). On non-Windows platforms the functions return `ErrUnsupportedPlatform`.
