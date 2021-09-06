package adoc

import (
	"encoding/xml"
	"fmt"
)

// Adoc is administrative divisions of China
// https://en.wikipedia.org/wiki/Administrative_divisions_of_China
type Adoc struct {
	XMLName xml.Name `json:"-" xml:"adoc"`
	Code    int64    `json:"code" xml:"code"`
	Parent  int64    `json:"parent" xml:"parent"`
	Name    string   `json:"name" xml:"name"`
}

func (a *Adoc) String() string {
	return fmt.Sprintf(`%d,%d,"%s"`, a.Code, a.Parent, a.Name)
}
