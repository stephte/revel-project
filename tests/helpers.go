package tests

import (
	"net/http"
	"strings"
	"bytes"
	"fmt"
)

func getRequest(typ string, urlEnd string, body []byte) *http.Request {
	typ = strings.ToUpper(typ)
	url := fmt.Sprintf("%s%s", urlBase(), urlEnd)

	req, e := http.NewRequest(typ, url, bytes.NewBuffer(body))
	if e != nil {
		panic(e)
	}

	if typ == "POST" || typ =="PUT" || typ == "PATCH" {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}