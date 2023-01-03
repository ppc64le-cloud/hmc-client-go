package hmc

import "golang.org/x/tools/blog/atom"

type (
	Feed struct {
		*atom.Feed
		Entry []*Entry `xml:"entry"`
	}

	Entry struct {
		*atom.Entry
		Content *Content `xml:"content"`
	}

	Content struct {
		atom.Text
		ManagedSystem ManagedSystem `xml:"ManagedSystem"`
	}

	ManagedSystem struct {
		ActivatedLevel string `xml:"ActivatedLevel"`
	}
)
