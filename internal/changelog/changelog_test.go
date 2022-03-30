package changelog

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewEntry(t *testing.T) {
	r := strings.NewReader(`---
channel: default
priority: 64
---
My changelog entry`)

	ee := &Entry{
		FrontMatter: FrontMatter{
			Channel:  "default",
			Priority: 64,
		},
		Text: "My changelog entry",
	}
	e, err := NewEntry(r)

	if err != nil || !reflect.DeepEqual(e, ee) {
		t.Errorf("entry must be %q, got %q", ee, e)
	}
}

func TestNewEntryError(t *testing.T) {
	r := strings.NewReader(`My changelog entry`)

	if _, err := NewEntry(r); err == nil {
		t.Errorf("a front matter must be present in the input data")
	}
}
