package gmail

import (
	"context"

	"google.golang.org/api/gmail/v1"
)

func (c *Client) CreateLabel(ctx context.Context, label string) error {
	c.labelmap.mu.Lock()
	defer c.labelmap.mu.Unlock()

	if _, ok := c.labelmap.name2id[label]; ok {
		return nil
	}

	newLabel, err := c.svc.Users.Labels.Create("me", &gmail.Label{Name: label}).Context(ctx).Do()
	if err != nil {
		return err
	}

	c.labelmap.id2name[newLabel.Id] = newLabel.Name
	c.labelmap.name2id[newLabel.Name] = newLabel.Id

	return nil
}
