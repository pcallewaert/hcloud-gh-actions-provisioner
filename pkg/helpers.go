package pkg

import (
	"strings"

	"github.com/google/go-github/v33/github"
)

func RunnerLabelsContains(labels []*github.RunnerLabels, e string) (bool, string) {
	for _, a := range labels {
		if strings.HasPrefix(a.GetName(), e) {
			return true, a.GetName()
		}
	}
	return false, ""
}
