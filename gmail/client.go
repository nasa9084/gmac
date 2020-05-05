package gmail

import (
	"context"
	"sync"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// newGmailService is replacable function for testing purpose.
var newGmailService = gmail.NewService

// Client is a gmail configuration client with label id <=> name map cache.
type Client struct {
	svc *gmail.Service

	// this client is assumet to be used in a command-line tool,
	// so this label map cache is only updated when client is
	// created, label is created or label is deleted.
	labelmap *labelmap
}

func New(ctx context.Context, oauthConfig *oauth2.Config, token *oauth2.Token) (*Client, error) {
	svc, err := newGmailService(ctx, option.WithTokenSource(oauthConfig.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}
	m, err := newLabelMap(ctx, svc)
	if err != nil {
		return nil, err
	}
	c := &Client{
		svc:      svc,
		labelmap: m,
	}
	return c, nil
}

type labelmap struct {
	mu sync.RWMutex

	id2name map[string]string
	name2id map[string]string
}

func newLabelMap(ctx context.Context, svc *gmail.Service) (*labelmap, error) {
	resp, err := svc.Users.Labels.List("me").Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	m := labelmap{
		id2name: map[string]string{},
		name2id: map[string]string{},
	}
	for _, label := range resp.Labels {
		m.id2name[label.Id] = label.Name
		m.name2id[label.Name] = label.Id
	}
	return &m, nil
}

func (m *labelmap) getNameByID(id string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if name, ok := m.id2name[id]; ok {
		return name
	}
	return ""
}

func (m *labelmap) getIDByName(name string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if id, ok := m.name2id[name]; ok {
		return id
	}
	return ""
}
