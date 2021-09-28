package models

import (
	"fmt"
)

// TunnelsResponse is exported, it models the data we receive.
type TunnelsResponse struct {
	Tunnels []Tunnel `json:"tunnels"`
}

//TextOutput is exported,it formats the data to plain text.
func (c TunnelsResponse) TextOutput() string {
	p := fmt.Sprintf(
		"Tunnels: %s\n",
		c.Tunnels)
	return p
}
