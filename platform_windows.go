//go:build windows

package credui

func newPlatformAPI(logger Logger) platformAPI {
	if logger != nil {
		logger.Debugf("newPlatformAPI: windows -> winAPIs")
	}
	return newWinAPIs(logger)
}
