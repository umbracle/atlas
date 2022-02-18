package userdata

import _ "embed"

//go:embed userdata.sh.tmpl
var userData string

func GetUserData() string {
	return userData
}
