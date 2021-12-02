package bip32ed25519

import "regexp"

var pathRegexStr = "^m(\\/[0-9]+')+$"

func validPath(path string) bool {
	regex, _ := regexp.Compile(pathRegexStr)
	match := regex.MatchString(path)
	return match
}
