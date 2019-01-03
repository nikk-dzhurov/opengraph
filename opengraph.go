package opengraph

import (
	"io"
	"net/url"
	"path"
	"path/filepath"

	"golang.org/x/net/html"
)

// OpenGraph represents web page information according to OGP <ogp.me>,
// and some more additional informations like URL.Host and so.
type OpenGraph struct {

	// Basics
	Title    string
	Type     string
	URL      URL
	SiteName string

	// Structures
	Image []*OGImage
	Video []*OGVideo
	Audio []*OGAudio

	// Optionals
	Description string
	Determiner  string // TODO: enum?
	Locale      string
	LocaleAlt   []string

	// Additionals
	Favicon string

	// Utils
	Error      error        `json:"-"`
}

// URL includes *url.URL
type URL struct {
	Source string
	*url.URL
}

// New creates new OpenGraph struct with specified URL.
func New(rawurl string) *OpenGraph {
	og := new(OpenGraph)
	og.Image = []*OGImage{}
	og.Video = []*OGVideo{}
	og.Audio = []*OGAudio{}
	og.LocaleAlt = []string{}
	og.Favicon = "/favicon.ico"
	u, err := url.Parse(rawurl)
	if err != nil {
		og.Error = err
		return og
	}
	og.URL = URL{Source: u.String(), URL: u}
	return og
}

// Parse parses http.Response.Body and construct OpenGraph informations.
// Caller should close body after it get parsed.
func (og *OpenGraph) Parse(body io.Reader) error {
	if og.Error != nil {
		return og.Error
	}
	node, err := html.Parse(body)
	if err != nil {
		return err
	}
	og.walk(node)
	return nil
}

func (og *OpenGraph) satisfied() bool {
	return false
}

func (og *OpenGraph) walk(n *html.Node) error {
	if og.satisfied() {
		return nil
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "body":
			return nil
		case "title":
			return TitleTag(n).Contribute(og)
		case "meta":
			return MetaTag(n).Contribute(og)
		case "link":
			return LinkTag(n).Contribute(og)
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		og.walk(child)
	}

	return nil
}

// ToAbsURL make og.Image and og.Favicon absolute URL if relative.
func (og *OpenGraph) ToAbsURL() *OpenGraph {
	for _, img := range og.Image {
		img.URL = og.abs(img.URL)
	}
	og.Favicon = og.abs(og.Favicon)
	return og
}

// abs make given URL absolute.
func (og *OpenGraph) abs(raw string) string {
	u, _ := url.Parse(raw)
	if u.IsAbs() {
		return raw
	}
	u.Scheme = og.URL.Scheme
	u.Host = og.URL.Host
	if !filepath.IsAbs(raw) {
		u.Path = path.Join(filepath.Dir(og.URL.Path), u.Path)
	}
	return u.String()
}

// Fulfill fulfills OG informations with some expectations.
func (og *OpenGraph) Fulfill() error {
	if og.SiteName == "" {
		og.SiteName = og.URL.Host
	}
	return nil
}
