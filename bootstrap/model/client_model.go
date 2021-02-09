package model

import (
	"github.com/ggoop/mdf/framework/md"
)

type Client struct {
	md.Model
	Name   string
	Secret string
	Memo   string
	Status string
}

func (t Client) TableName() string {
	return "sys_clients"
}
func (s *Client) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".client", Domain: MD_DOMAIN, Name: "端"}
}
