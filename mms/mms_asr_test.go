package mms

import (
	"context"
	"testing"
)

func TestMMSASR_ProcessFiles(t *testing.T) {

}

func TestMMSASR_checkLanguagea(t *testing.T) {
	checkLanguageTest("fra", "en", "en", t)
	checkLanguageTest("fra", "", "fra", t)
	checkLanguageTest("npi", "", "npi", t)
	checkLanguageTest("abc", "", "pam", t)
}

func checkLanguageTest(lang string, sttLang string, expectLang string, t *testing.T) {
	ctx := context.Background()
	resultLang, status := checkLanguage(ctx, lang, sttLang)
	if status.IsErr {
		t.Error(status)
	}
	if resultLang != expectLang {
		t.Error("actual:", resultLang, "expected:", expectLang)
	}
}
