package cmd

import "testing"

func TestRootHelp_DefaultsToTraditionalChinese(t *testing.T) {
	output := runRootHelp(t, "")
	assertGolden(t, "help/root_zh_tw.golden", output)
}

func TestRootHelp_SupportsEnglishViaConfig(t *testing.T) {
	output := runRootHelp(t, `version = 1

[i18n]
language = "en"
`)
	assertGolden(t, "help/root_en.golden", output)
}
