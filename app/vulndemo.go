package main

import "golang.org/x/text/language"

// Temporary: reachable call into a vulnerable x/text (GO-2021-0113) to demo the govulncheck gate.
func init() { _, _, _ = language.ParseAcceptLanguage("en-US") }
