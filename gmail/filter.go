package gmail

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/api/gmail/v1"
)

const ResourceTypeFilter = "Filter"

type FilterResource struct {
	Filters []Filter `yaml:"filters"`
}

type Filter struct {
	id string `yaml:"-"`

	Criteria FilterCriteria `yaml:"criteria"`
	Action   FilterAction   `yaml:"action"`
}

type FilterCriteria struct {
	From          string `yaml:"from,omitempty"`
	To            string `yaml:"to,omitempty"`
	Subject       string `yaml:"subject,omitempty"`
	Query         string `yaml:"query,omitempty"`
	NegatedQuery  string `yaml:"negated_query,omitempty"`
	LargerThan    int64  `yaml:"larger_than,omitempty"`
	SmallerThan   int64  `yaml:"smaller_than,omitempty"`
	HasAttachment bool   `yaml:"has_attachment,omitempty"`
	ExcludeChats  bool   `yaml:"exclude_chats,omitempty"`
}

type FilterAction struct {
	Archive         bool                  `yaml:"archive,omitempty"`
	MarkAsRead      bool                  `yaml:"mark_as_read,omitempty"`
	Star            bool                  `yaml:"star,omitempty"`
	AddLabel        string                `yaml:"add_label,omitempty"`
	ForwardTo       string                `yaml:"forward_to,omitempty"`
	Delete          bool                  `yaml:"delete,omitempty"`
	NeverMarkAsSpam bool                  `yaml:"never_mark_as_spam,omitempty"`
	Important       FilterActionImportant `yaml:"important,omitempty"`
	Category        string                `yaml:"category,omitempty"`
}

type FilterActionImportant string

const (
	FilterActionImportantAlways FilterActionImportant = "always"
	FilterActionImportantNever  FilterActionImportant = "never"
)

func (c *Client) ListFilters(ctx context.Context) ([]Filter, error) {
	resp, err := c.svc.Users.Settings.Filters.List("me").Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	var filters []Filter
	for _, gf := range resp.Filter {
		filters = append(filters, c.convertFilterFromGmail(gf))
	}
	return filters, nil
}

func (c *Client) CreateFilter(ctx context.Context, filter Filter) error {
	if filter.Action.AddLabel != "" {
		if err := c.CreateLabel(ctx, filter.Action.AddLabel); err != nil {
			return err
		}
	}

	gf, err := c.convertFilterToGmail(filter)
	if err != nil {
		return err
	}
	if _, err := c.svc.Users.Settings.Filters.Create("me", gf).Context(ctx).Do(); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteAllFilter(ctx context.Context) error {
	filters, err := c.ListFilters(ctx)
	if err != nil {
		return err
	}

	for _, filter := range filters {
		if err := c.DeleteFilterByID(ctx, filter.id); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) DeleteFilterByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id must be non-empty")
	}
	if err := c.svc.Users.Settings.Filters.Delete("me", id).Context(ctx).Do(); err != nil {
		return err

	}
	return nil
}

func (c *Client) convertFilterFromGmail(gf *gmail.Filter) Filter {
	var f Filter
	f.id = gf.Id
	// criteria
	f.Criteria = FilterCriteria{
		From:          gf.Criteria.From,
		To:            gf.Criteria.To,
		Subject:       gf.Criteria.Subject,
		Query:         gf.Criteria.Query,
		NegatedQuery:  gf.Criteria.NegatedQuery,
		HasAttachment: gf.Criteria.HasAttachment,
		ExcludeChats:  gf.Criteria.ExcludeChats,
	}
	if gf.Criteria.SizeComparison == "larger" {
		f.Criteria.LargerThan = gf.Criteria.Size
	} else if gf.Criteria.SizeComparison == "smaller" {
		f.Criteria.SmallerThan = gf.Criteria.Size
	}
	// actions
	f.Action.ForwardTo = gf.Action.Forward
	for _, id := range gf.Action.AddLabelIds {
		switch id {
		case "TRASH":
			f.Action.Delete = true
		case "IMPORTANT":
			f.Action.Important = FilterActionImportantAlways
		case "STARRED":
			f.Action.Star = true
		case "CATEGORY_PERSONAL":
			f.Action.Category = "primary"
		case "CATEGORY_SOCIAL":
			f.Action.Category = "social"
		case "CATEGORY_UPDATES":
			f.Action.Category = "updates"
		case "CATEGORY_FORUMS":
			f.Action.Category = "forums"
		case "CATEGORY_PROMOTIONS":
			f.Action.Category = "promotions"
		default:
			f.Action.AddLabel = c.labelmap.getNameByID(id)
		}
	}
	for _, id := range gf.Action.RemoveLabelIds {
		switch id {
		case "INBOX":
			f.Action.Archive = true
		case "UNREAD":
			f.Action.MarkAsRead = true
		case "SPAM":
			f.Action.NeverMarkAsSpam = true
		case "IMPORTANT":
			f.Action.Important = FilterActionImportantNever
		}
	}
	return f
}

func (c *Client) convertFilterToGmail(filter Filter) (*gmail.Filter, error) {
	gf := &gmail.Filter{
		Criteria: &gmail.FilterCriteria{
			From:          filter.Criteria.From,
			To:            filter.Criteria.To,
			Subject:       filter.Criteria.Subject,
			Query:         filter.Criteria.Query,
			NegatedQuery:  filter.Criteria.NegatedQuery,
			HasAttachment: filter.Criteria.HasAttachment,
			ExcludeChats:  filter.Criteria.ExcludeChats,
		},
		Action: &gmail.FilterAction{},
	}
	if filter.Criteria.LargerThan > 0 {
		gf.Criteria.SizeComparison = "larger"
		gf.Criteria.Size = filter.Criteria.LargerThan
	} else if filter.Criteria.SmallerThan > 0 {
		gf.Criteria.SizeComparison = "smaller"
		gf.Criteria.Size = filter.Criteria.SmallerThan
	}
	gf.Action.Forward = filter.Action.ForwardTo
	if filter.Action.AddLabel != "" {
		gf.Action.AddLabelIds = append(gf.Action.AddLabelIds, c.labelmap.getIDByName(filter.Action.AddLabel))
	}
	if filter.Action.Delete {
		gf.Action.AddLabelIds = append(gf.Action.AddLabelIds, "TRASH")
	}
	switch filter.Action.Important {
	case "": // nothing to do
	case FilterActionImportantAlways:
		gf.Action.AddLabelIds = append(gf.Action.AddLabelIds, "IMPORTANT")
	case FilterActionImportantNever:
		gf.Action.RemoveLabelIds = append(gf.Action.RemoveLabelIds, "IMPORTANT")
	default:
		return nil, fmt.Errorf("unknown action.important value: %s", filter.Action.Important)
	}
	if filter.Action.Star {
		gf.Action.AddLabelIds = append(gf.Action.AddLabelIds, "STARRED")
	}
	switch filter.Action.Category {
	case "": // nothing to do
	case "primary", "personal", "mail":
		gf.Action.AddLabelIds = append(gf.Action.AddLabelIds, "CATEGORY_PERSONAL")
	case "social":
		gf.Action.AddLabelIds = append(gf.Action.AddLabelIds, "CATEGORY_SOCIAL")
	case "update", "updates", "new":
		gf.Action.AddLabelIds = append(gf.Action.AddLabelIds, "CATEGORY_UPDATES")
	case "forum", "forums":
		gf.Action.AddLabelIds = append(gf.Action.AddLabelIds, "CATEGORY_FORUMS")
	case "promotion", "promotions":
		gf.Action.AddLabelIds = append(gf.Action.AddLabelIds, "CATEGORY_PROMOTIONS")
	default:
		return nil, fmt.Errorf("unknown action.category value: %s", filter.Action.Category)
	}
	if filter.Action.Archive {
		gf.Action.RemoveLabelIds = append(gf.Action.RemoveLabelIds, "INBOX")
	}
	if filter.Action.MarkAsRead {
		gf.Action.RemoveLabelIds = append(gf.Action.RemoveLabelIds, "UNREAD")
	}
	if filter.Action.NeverMarkAsSpam {
		gf.Action.RemoveLabelIds = append(gf.Action.RemoveLabelIds, "SPAM")
	}
	return gf, nil
}

func (f Filter) String() string {
	return f.Criteria.String() + " => " + f.Action.String()
}

var nonParen = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (criteria FilterCriteria) String() string {
	var cq []string
	if criteria.From != "" {
		if nonParen.MatchString(criteria.From) {
			cq = append(cq, "from:"+criteria.From)
		} else {
			cq = append(cq, fmt.Sprintf("from:(%s)", criteria.From))
		}
	}
	if criteria.To != "" {
		if nonParen.MatchString(criteria.To) {
			cq = append(cq, "to:"+criteria.To)
		} else {
			cq = append(cq, fmt.Sprintf("to:(%s)", criteria.To))
		}
	}
	if criteria.Subject != "" {
		if nonParen.MatchString(criteria.Subject) {
			cq = append(cq, "subject:"+criteria.Subject)
		} else {
			cq = append(cq, fmt.Sprintf("subject:(%s)", criteria.Subject))
		}
	}
	if criteria.Query != "" {
		cq = append(cq, criteria.Query)
	}
	if criteria.NegatedQuery != "" {
		if nonParen.MatchString(criteria.NegatedQuery) {
			cq = append(cq, "-"+criteria.NegatedQuery)
		} else {
			cq = append(cq, fmt.Sprintf("-{%s}", criteria.NegatedQuery))
		}
	}
	if criteria.HasAttachment {
		cq = append(cq, "has:attachment")
	}
	if criteria.ExcludeChats {
		cq = append(cq, "-in:chats")
	}
	if criteria.LargerThan > 0 {
		cq = append(cq, "larger:"+strconv.FormatInt(criteria.LargerThan, 10))
	} else if criteria.SmallerThan > 0 {
		cq = append(cq, "smaller:"+strconv.FormatInt(criteria.SmallerThan, 10))
	}
	return strings.Join(cq, " ")
}

func (action FilterAction) String() string {
	var aq []string

	if action.Archive {
		aq = append(aq, "Skip Inbox")
	}

	if action.MarkAsRead {
		aq = append(aq, "Mark as read")
	}

	if action.Star {
		aq = append(aq, "Star it")
	}

	if action.AddLabel != "" {
		aq = append(aq, fmt.Sprintf(`Apply label "%s"`, action.AddLabel))
	}

	if action.ForwardTo != "" {
		aq = append(aq, "Forward to "+action.ForwardTo)
	}

	if action.Delete {
		aq = append(aq, "Delete it")
	}

	if action.NeverMarkAsSpam {
		aq = append(aq, "Never send it to Spam")
	}

	switch action.Important {
	case FilterActionImportantAlways:
		aq = append(aq, "Mark it as important")
	case FilterActionImportantNever:
		aq = append(aq, "Never mark it as important")
	}

	if action.Category != "" {
		aq = append(aq, "Categorize as "+action.Category)
	}

	return strings.Join(aq, ", ")
}
