package gmail

import (
	"reflect"
	"sort"
	"strconv"
	"testing"

	"google.golang.org/api/gmail/v1"
)

func TestConvertFilterFromGmail(t *testing.T) {
	tests := []struct {
		label string
		input *gmail.Filter
		want  Filter
	}{
		{
			label: "check string-based criteria",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{
					From:         "fromAddrFoo",
					To:           "toAddrFoo",
					Subject:      "subjectStringFoo",
					Query:        "queryStringFoo",
					NegatedQuery: "negatedQueryStringFoo",
				},
				Action: &gmail.FilterAction{},
			},
			want: Filter{
				Criteria: FilterCriteria{
					From:         "fromAddrFoo",
					To:           "toAddrFoo",
					Subject:      "subjectStringFoo",
					Query:        "queryStringFoo",
					NegatedQuery: "negatedQueryStringFoo",
				},
			},
		},
		{
			label: "check smaller-than criteria",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{
					Size:           2000000,
					SizeComparison: "smaller",
				},
				Action: &gmail.FilterAction{},
			},
			want: Filter{
				Criteria: FilterCriteria{
					SmallerThan: 2000000,
				},
			},
		},
		{
			label: "check larger-than criteria",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{
					Size:           4000000,
					SizeComparison: "larger",
				},
				Action: &gmail.FilterAction{},
			},
			want: Filter{
				Criteria: FilterCriteria{
					LargerThan: 4000000,
				},
			},
		},
		{
			label: "check bool-based criteria",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{
					HasAttachment: true,
					ExcludeChats:  true,
				},
				Action: &gmail.FilterAction{},
			},
			want: Filter{
				Criteria: FilterCriteria{
					HasAttachment: true,
					ExcludeChats:  true,
				},
			},
		},
		{
			label: "check archive action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					RemoveLabelIds: []string{"INBOX"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Archive: true,
				},
			},
		},
		{
			label: "check mark as read action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					RemoveLabelIds: []string{"UNREAD"},
				},
			},
			want: Filter{
				Action: FilterAction{
					MarkAsRead: true,
				},
			},
		},
		{
			label: "check star action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"STARRED"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Star: true,
				},
			},
		},
		{
			label: "check add label action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"Label_10"},
				},
			},
			want: Filter{
				Action: FilterAction{
					AddLabel: "LabelFoo",
				},
			},
		},
		{
			label: "check forward to action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					Forward: "forwardTo@example.com",
				},
			},
			want: Filter{
				Action: FilterAction{
					ForwardTo: "forwardTo@example.com",
				},
			},
		},
		{
			label: "check delete action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"TRASH"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Delete: true,
				},
			},
		},
		{
			label: "check never mark as spam action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					RemoveLabelIds: []string{"SPAM"},
				},
			},
			want: Filter{
				Action: FilterAction{
					NeverMarkAsSpam: true,
				},
			},
		},
		{
			label: "check always important action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"IMPORTANT"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Important: FilterActionImportantAlways,
				},
			},
		},
		{
			label: "check never mark as important action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					RemoveLabelIds: []string{"IMPORTANT"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Important: FilterActionImportantNever,
				},
			},
		},
		{
			label: "check categorize (primary) action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_PERSONAL"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Category: "primary",
				},
			},
		},
		{
			label: "check categorize (social) action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_SOCIAL"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Category: "social",
				},
			},
		},
		{
			label: "check categorize (updates) action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_UPDATES"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Category: "updates",
				},
			},
		},
		{
			label: "check categorize (forums) action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_FORUMS"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Category: "forums",
				},
			},
		},
		{
			label: "check categorize (promotions) action",
			input: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_PROMOTIONS"},
				},
			},
			want: Filter{
				Action: FilterAction{
					Category: "promotions",
				},
			},
		},
	}

	c := &Client{
		labelmap: &labelmap{
			id2name: map[string]string{
				"Label_10": "LabelFoo",
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"."+tt.label, func(t *testing.T) {
			got := c.convertFilterFromGmail(tt.input)
			// check Criteria
			if got.Criteria.From != tt.want.Criteria.From {
				t.Errorf("unexpected Criteria.From: %s != %s",
					got.Criteria.From,
					tt.want.Criteria.From,
				)
				return
			}
			if got.Criteria.To != tt.want.Criteria.To {
				t.Errorf("unexpected Criteria.To: %s != %s",
					got.Criteria.To,
					tt.want.Criteria.To,
				)
				return
			}
			if got.Criteria.Subject != tt.want.Criteria.Subject {
				t.Errorf("unexpected Criteria.Subject: %s != %s",
					got.Criteria.Subject,
					tt.want.Criteria.Subject,
				)
				return
			}
			if got.Criteria.Query != tt.want.Criteria.Query {
				t.Errorf("unexpected Criteria.Query: %s != %s",
					got.Criteria.Query,
					tt.want.Criteria.Query,
				)
				return
			}
			if got.Criteria.NegatedQuery != tt.want.Criteria.NegatedQuery {
				t.Errorf("unexpected Criteria.NegatedQuery: %s != %s",
					got.Criteria.NegatedQuery,
					tt.want.Criteria.NegatedQuery,
				)
				return
			}
			if got.Criteria.LargerThan != tt.want.Criteria.LargerThan {
				t.Errorf("unexpected Criteria.LagerThan: %d != %d",
					got.Criteria.LargerThan,
					tt.want.Criteria.LargerThan,
				)
				return
			}
			if got.Criteria.SmallerThan != tt.want.Criteria.SmallerThan {
				t.Errorf("unexpected Criteria.SmallerThan: %d != %d",
					got.Criteria.SmallerThan,
					tt.want.Criteria.SmallerThan,
				)
				return
			}
			if got.Criteria.HasAttachment != tt.want.Criteria.HasAttachment {
				t.Errorf("unexpected Criteria.HasAttachment: %t != %t",
					got.Criteria.HasAttachment,
					tt.want.Criteria.HasAttachment,
				)
				return
			}
			if got.Criteria.ExcludeChats != tt.want.Criteria.ExcludeChats {
				t.Errorf("unexpected Criteria.ExcludeChats: %t != %t",
					got.Criteria.ExcludeChats,
					tt.want.Criteria.ExcludeChats,
				)
				return
			}

			// check Action
			if got.Action.Archive != tt.want.Action.Archive {
				t.Errorf("unexpected Action.Archive: %t != %t",
					got.Action.Archive,
					tt.want.Action.Archive,
				)
				return
			}
			if got.Action.MarkAsRead != tt.want.Action.MarkAsRead {
				t.Errorf("unexpected Action.MarkAsRead: %t != %t",
					got.Action.MarkAsRead,
					tt.want.Action.MarkAsRead,
				)
				return
			}
			if got.Action.Star != tt.want.Action.Star {
				t.Errorf("unexpected Action.Star: %t != %t",
					got.Action.Star,
					tt.want.Action.Star,
				)
				return
			}
			if got.Action.AddLabel != tt.want.Action.AddLabel {
				t.Errorf("unexpected Action.AddLabel: %s != %s",
					got.Action.AddLabel,
					tt.want.Action.AddLabel,
				)
				return
			}
			if got.Action.ForwardTo != tt.want.Action.ForwardTo {
				t.Errorf("unexpected Action.ForwardTo: %s != %s",
					got.Action.ForwardTo,
					tt.want.Action.ForwardTo,
				)
				return
			}
			if got.Action.Delete != tt.want.Action.Delete {
				t.Errorf("unexpected Action.Delete: %t != %t",
					got.Action.Delete,
					tt.want.Action.Delete,
				)
				return
			}
			if got.Action.NeverMarkAsSpam != tt.want.Action.NeverMarkAsSpam {
				t.Errorf("unexpected Action.NeverMarkAsSpam: %t != %t",
					got.Action.NeverMarkAsSpam,
					tt.want.Action.NeverMarkAsSpam,
				)
				return
			}
			if got.Action.Important != tt.want.Action.Important {
				t.Errorf("unexpected Action.Important: %s != %s",
					got.Action.Important,
					tt.want.Action.Important,
				)
				return
			}
			if got.Action.Category != tt.want.Action.Category {
				t.Errorf("unexpected Action.Category: %s != %s",
					got.Action.Category,
					tt.want.Action.Category,
				)
				return
			}
		})
	}
}

func TestConvertFilterToGmail(t *testing.T) {
	tests := []struct {
		label string
		input Filter
		want  *gmail.Filter
	}{
		{
			label: "check string-based criteria",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{
					From:         "fromAddrFoo",
					To:           "toAddrFoo",
					Subject:      "subjectStringFoo",
					Query:        "queryStringFoo",
					NegatedQuery: "negatedQueryStringFoo",
				},
				Action: &gmail.FilterAction{},
			},
			input: Filter{
				Criteria: FilterCriteria{
					From:         "fromAddrFoo",
					To:           "toAddrFoo",
					Subject:      "subjectStringFoo",
					Query:        "queryStringFoo",
					NegatedQuery: "negatedQueryStringFoo",
				},
			},
		},
		{
			label: "check smaller-than criteria",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{
					Size:           2000000,
					SizeComparison: "smaller",
				},
				Action: &gmail.FilterAction{},
			},
			input: Filter{
				Criteria: FilterCriteria{
					SmallerThan: 2000000,
				},
			},
		},
		{
			label: "check larger-than criteria",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{
					Size:           4000000,
					SizeComparison: "larger",
				},
				Action: &gmail.FilterAction{},
			},
			input: Filter{
				Criteria: FilterCriteria{
					LargerThan: 4000000,
				},
			},
		},
		{
			label: "check bool-based criteria",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{
					HasAttachment: true,
					ExcludeChats:  true,
				},
				Action: &gmail.FilterAction{},
			},
			input: Filter{
				Criteria: FilterCriteria{
					HasAttachment: true,
					ExcludeChats:  true,
				},
			},
		},
		{
			label: "check archive action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					RemoveLabelIds: []string{"INBOX"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Archive: true,
				},
			},
		},
		{
			label: "check mark as read action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					RemoveLabelIds: []string{"UNREAD"},
				},
			},
			input: Filter{
				Action: FilterAction{
					MarkAsRead: true,
				},
			},
		},
		{
			label: "check star action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"STARRED"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Star: true,
				},
			},
		},
		{
			label: "check add label action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"Label_10"},
				},
			},
			input: Filter{
				Action: FilterAction{
					AddLabel: "LabelFoo",
				},
			},
		},
		{
			label: "check forward to action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					Forward: "forwardTo@example.com",
				},
			},
			input: Filter{
				Action: FilterAction{
					ForwardTo: "forwardTo@example.com",
				},
			},
		},
		{
			label: "check delete action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"TRASH"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Delete: true,
				},
			},
		},
		{
			label: "check never mark as spam action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					RemoveLabelIds: []string{"SPAM"},
				},
			},
			input: Filter{
				Action: FilterAction{
					NeverMarkAsSpam: true,
				},
			},
		},
		{
			label: "check always important action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"IMPORTANT"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Important: FilterActionImportantAlways,
				},
			},
		},
		{
			label: "check never mark as important action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					RemoveLabelIds: []string{"IMPORTANT"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Important: FilterActionImportantNever,
				},
			},
		},
		{
			label: "check categorize (primary) action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_PERSONAL"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Category: "primary",
				},
			},
		},
		{
			label: "check categorize (social) action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_SOCIAL"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Category: "social",
				},
			},
		},
		{
			label: "check categorize (updates) action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_UPDATES"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Category: "updates",
				},
			},
		},
		{
			label: "check categorize (forums) action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_FORUMS"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Category: "forums",
				},
			},
		},
		{
			label: "check categorize (promotions) action",
			want: &gmail.Filter{
				Criteria: &gmail.FilterCriteria{},
				Action: &gmail.FilterAction{
					AddLabelIds: []string{"CATEGORY_PROMOTIONS"},
				},
			},
			input: Filter{
				Action: FilterAction{
					Category: "promotions",
				},
			},
		},
	}

	c := &Client{
		labelmap: &labelmap{
			name2id: map[string]string{
				"LabelFoo": "Label_10",
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"."+tt.label, func(t *testing.T) {
			got, err := c.convertFilterToGmail(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			// check Criteria
			if got.Criteria.From != tt.want.Criteria.From {
				t.Errorf("unexpected Criteria.From: %s != %s",
					got.Criteria.From,
					tt.want.Criteria.From,
				)
				return
			}
			if got.Criteria.To != tt.want.Criteria.To {
				t.Errorf("unexpected Criteria.To: %s != %s",
					got.Criteria.To,
					tt.want.Criteria.To,
				)
				return
			}
			if got.Criteria.Subject != tt.want.Criteria.Subject {
				t.Errorf("unexpected Criteria.Subject: %s != %s",
					got.Criteria.Subject,
					tt.want.Criteria.Subject,
				)
				return
			}
			if got.Criteria.Query != tt.want.Criteria.Query {
				t.Errorf("unexpected Criteria.Query: %s != %s",
					got.Criteria.Query,
					tt.want.Criteria.Query,
				)
				return
			}
			if got.Criteria.NegatedQuery != tt.want.Criteria.NegatedQuery {
				t.Errorf("unexpected Criteria.NegatedQuery: %s != %s",
					got.Criteria.NegatedQuery,
					tt.want.Criteria.NegatedQuery,
				)
				return
			}
			if got.Criteria.HasAttachment != tt.want.Criteria.HasAttachment {
				t.Errorf("unexpected Criteria.HasAttachment: %t != %t",
					got.Criteria.HasAttachment,
					tt.want.Criteria.HasAttachment,
				)
				return
			}
			if got.Criteria.ExcludeChats != tt.want.Criteria.ExcludeChats {
				t.Errorf("unexpected Criteria.ExcludeChats: %t != %t",
					got.Criteria.ExcludeChats,
					tt.want.Criteria.ExcludeChats,
				)
				return
			}
			if got.Criteria.SizeComparison != tt.want.Criteria.SizeComparison {
				t.Errorf("unexpected Criteria.SizeComparison: %s != %s",
					got.Criteria.SizeComparison,
					tt.want.Criteria.SizeComparison,
				)
				return
			}
			if got.Criteria.Size != tt.want.Criteria.Size {
				t.Errorf("unexpected Criteria.Size: %d != %d",
					got.Criteria.Size,
					tt.want.Criteria.Size,
				)
				return
			}

			// check Action
			sort.Strings(got.Action.AddLabelIds)
			sort.Strings(tt.want.Action.AddLabelIds)
			if !reflect.DeepEqual(got.Action.AddLabelIds, tt.want.Action.AddLabelIds) {
				t.Errorf("unexpected Action.AddLabelIds:\n  got:  %v\n  want: %v",
					got.Action.AddLabelIds,
					tt.want.Action.AddLabelIds,
				)
				return
			}
			sort.Strings(got.Action.RemoveLabelIds)
			sort.Strings(tt.want.Action.RemoveLabelIds)
			if !reflect.DeepEqual(got.Action.RemoveLabelIds, tt.want.Action.RemoveLabelIds) {
				t.Errorf("unexpected Action.RemoveLabelIds:\n  got:  %v\n  want: %v",
					got.Action.RemoveLabelIds,
					tt.want.Action.RemoveLabelIds,
				)
				return
			}
		})
	}
}
