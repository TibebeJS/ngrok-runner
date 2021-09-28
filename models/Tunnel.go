package models

import (
	"fmt"
)

type Tunnel struct {
	Protocol  string `json:"proto"`
	PublicURL string `json:"public_url"`
}

func (c Tunnel) TextOutput() string {
	p := fmt.Sprintf(
		"Protocol: %s\nURL : %s\n",
		c.Protocol, c.PublicURL)
	return p
}
