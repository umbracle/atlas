package userdata

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed userdata.sh.tmpl
var userData string

func GetUserData(downloadPath string) (string, error) {
	tmpl, err := template.New("test").Parse(userData)
	if err != nil {
		return "", err
	}

	config := map[string]interface{}{
		"DownloadLink": downloadPath,
	}

	out := bytes.NewBuffer(nil)
	if err := tmpl.Execute(out, config); err != nil {
		return "", err
	}
	return out.String(), nil
}
