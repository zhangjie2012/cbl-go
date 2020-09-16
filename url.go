package cbl

import (
	"net/url"
	"path"
)

// UrlJoin via https://stackoverflow.com/questions/34668012/combine-url-paths-with-path-join
func UrlJoin(host string, subs ...string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	t := []string{u.Path}
	t = append(t, subs...)
	u.Path = path.Join(t...)
	return u.String(), nil
}
