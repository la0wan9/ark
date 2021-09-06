package adoc

import adocv1 "github.com/la0wan9/ark/pkg/adoc/v1"

// FromAdocToMessage transforms format from Adoc to Message
func FromAdocToMessage(a *Adoc) *adocv1.Adoc {
	return &adocv1.Adoc{
		Code:   a.Code,
		Parent: a.Parent,
		Name:   a.Name,
	}
}

// FromMessageToAdoc transforms format from Message to Adoc
func FromMessageToAdoc(a *adocv1.Adoc) *Adoc {
	return &Adoc{
		Code:   a.Code,
		Parent: a.Parent,
		Name:   a.Name,
	}
}
