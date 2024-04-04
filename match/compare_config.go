package match

import "golang.org/x/text/unicode/norm"

const (
	nothing = iota
	normalize
	remove
)

type config struct {
	lowerCase         bool
	removePromptChars bool
	removePunctuation bool
	doubleQuotes      int
	apostrophe        int
	hyphen            int
	diacritical       int       // if remove, must do NFD or NFKD
	normalizeType     norm.Form // NFC, NFD, NFKC, NFKD
}

func getConfig() config {
	var cfg config
	cfg.lowerCase = true
	cfg.removePromptChars = true
	cfg.removePunctuation = true
	cfg.doubleQuotes = remove
	cfg.apostrophe = remove
	cfg.hyphen = remove
	cfg.diacritical = normalize
	cfg.normalizeType = norm.NFD
	return cfg
}
