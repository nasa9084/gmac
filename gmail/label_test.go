package gmail

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const labelListResponseBody = `{
 "labels": [
  {
   "id": "INBOX",
   "name": "INBOX",
   "messageListVisibility": "hide",
   "labelListVisibility": "labelShow",
   "type": "system"
  },
  {
   "id": "IMPORTANT",
   "name": "IMPORTANT",
   "messageListVisibility": "hide",
   "labelListVisibility": "labelHide",
   "type": "system"
  },
  {
   "id": "UNREAD",
   "name": "UNREAD",
   "type": "system"
  },
  {
   "id": "Label_10",
   "name": "Foo",
   "messageListVisibility": "show",
   "labelListVisibility": "labelShow",
   "type": "user"
  },
  {
   "id": "Label_11",
   "name": "Bar",
   "type": "user"
  }
 ]
}`

func TestCreateLabel(t *testing.T) {
	const newLabelName = "newLabel"
	const newLabelID = "newLabelID"

	oauthSrv, oauthCfg, oauthToken := testOAuth(t)
	defer oauthSrv.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(labelListResponseBody))
		case http.MethodPost:
			var body struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.Name != newLabelName {
				t.Errorf("new label name in create request is unexpected: %s != %s", body.Name, newLabelName)
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			}{
				Name: newLabelName,
				ID:   newLabelID,
			}); err != nil {
				t.Fatal(err)
			}
		}
	}))
	defer srv.Close()

	newGmailService = func(ctx context.Context, opts ...option.ClientOption) (*gmail.Service, error) {
		opts = append(opts, option.WithEndpoint(srv.URL))
		return gmail.NewService(ctx, opts...)
	}

	ctx := context.Background()
	c, err := New(ctx, oauthCfg, oauthToken)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.CreateLabel(ctx, newLabelName); err != nil {
		t.Fatal(err)
	}

	wantLabelMap := labelmap{
		id2name: map[string]string{
			"INBOX":     "INBOX",
			"IMPORTANT": "IMPORTANT",
			"UNREAD":    "UNREAD",
			"Label_10":  "Foo",
			"Label_11":  "Bar",
			newLabelID:  newLabelName,
		},
		name2id: map[string]string{
			"INBOX":      "INBOX",
			"IMPORTANT":  "IMPORTANT",
			"UNREAD":     "UNREAD",
			"Foo":        "Label_10",
			"Bar":        "Label_11",
			newLabelName: newLabelID,
		},
	}
	if !reflect.DeepEqual(c.labelmap.id2name, wantLabelMap.id2name) {
		t.Errorf("unexpected labelmap.id2name:\n  got:  %#v\n  want: %#v", c.labelmap.id2name, wantLabelMap.id2name)
		return
	}
	if !reflect.DeepEqual(c.labelmap.name2id, wantLabelMap.name2id) {
		t.Errorf("unexpected labelmap.name2id:\n  got:  %#v\n  want: %#v", c.labelmap.name2id, wantLabelMap.name2id)
		return
	}
}
