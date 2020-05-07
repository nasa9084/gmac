package encoder

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"text/tabwriter"

	"github.com/goccy/go-yaml"
	"github.com/nasa9084/gmac/gmail"
)

var nonParen = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// FilterEncoder is an interface which encodes Filter object into string.
type FilterEncoder interface {
	Encode([]gmail.Filter) error
}

func NewFilterEncoder(w io.Writer, format string) FilterEncoder {
	switch format {
	case "yaml":
		return &yamlFilterEncoder{
			enc: yaml.NewEncoder(w),
		}
	case "wide":
		return &defaultFilterEncoder{
			w:      w,
			isWide: true,
		}
	default:
		return &defaultFilterEncoder{
			w: w,
		}
	}
}

// defaultFilterEncoder encodes Filter object into string by
// FilterCriteria.String() and FilterAction.String() methods.
// the criteria and action are abbreviated if isWide is false.
type defaultFilterEncoder struct {
	w io.Writer

	isWide bool
}

func (e *defaultFilterEncoder) Encode(filters []gmail.Filter) error {
	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 2, 1, ' ', 0)

	fmt.Fprint(w, "MATCHES\tACTION\n")

	for _, filter := range filters {
		fmt.Fprintf(w, "%s\t%s\n", e.encodeCriteria(filter.Criteria), e.encodeAction(filter.Action))

	}

	if err := w.Flush(); err != nil {
		return err
	}

	_, err := buf.WriteTo(e.w)
	return err
}

func (e *defaultFilterEncoder) encodeCriteria(criteria gmail.FilterCriteria) string {
	if e.isWide {
		return criteria.String()
	}
	return abbr(criteria.String(), 40)
}

func (e *defaultFilterEncoder) encodeAction(action gmail.FilterAction) string {
	if e.isWide {
		return action.String()
	}
	return abbr(action.String(), 40)
}

// yamlFilterEncoder encodes Filter object by marshaling to YAML format.
type yamlFilterEncoder struct {
	enc *yaml.Encoder
}

func (e *yamlFilterEncoder) Encode(filters []gmail.Filter) error {
	return e.enc.Encode(struct {
		Kind    string         `yaml:"kind"`
		Filters []gmail.Filter `yaml:"filters"`
	}{
		Kind:    "Filter",
		Filters: filters,
	})
}

func abbr(s string, maxLength int) string {
	if len(s) < maxLength {
		return s
	}
	return s[:maxLength-4] + "..."
}
