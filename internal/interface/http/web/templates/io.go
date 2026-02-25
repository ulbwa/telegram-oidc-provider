package templates

import _ "embed"

//go:embed login.html
var loginTemplate string

//go:embed error.html
var errorTemplate string

//go:embed consent.html
var consentTemplate string

func LoginTemplate() string {
	return loginTemplate
}

func ErrorTemplate() string {
	return errorTemplate
}

func ConsentTemplate() string {
	return consentTemplate
}
