package opengraph

import (
	"bytes"
	"io/ioutil"
	"path"
	"runtime"
	"testing"
)

func TestFetch(t *testing.T) {

	URL := "http://test.test/html/01.html"

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Error("Unexpected error")
		t.Fail()
		return
	}

	html := path.Join(path.Dir(filename), "./html/01.html")

	body, err := ioutil.ReadFile(html)
	if err != nil {
		t.Error("File not found", err)
		t.Fail()
		return
	}

	og := New(URL)
	err = og.Parse(bytes.NewReader(body))
	if err != nil {
		t.Error("Unexpected error", err)
		t.Fail()
		return
	}

	if og.Title != "Open Graph Title" {
		t.Error("Invalid Title")
		t.Fail()
		return
	}

	if og.Type != "website" {
		t.Error("Invalid Type")
		t.Fail()
		return
	}

	if og.URL.Source != URL {
		t.Error("Invalid URL")
		t.Fail()
		return
	}

	if len(og.Image) != 1 || og.Image[0].URL != "/images/01.png" {
		t.Error("Invalid Image")
		t.Fail()
		return
	}

	if og.Favicon != "/images/01.favicon.png" {
		t.Error("Invalid Favicon")
		t.Fail()
		return
	}
}
