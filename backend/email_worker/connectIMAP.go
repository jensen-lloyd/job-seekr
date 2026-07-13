package main

import (
	"crypto/tls"

	"github.com/emersion/go-imap/client"
)

func connectIMAP(server, username, password string) (*client.Client, error) {
	c, err := client.Dial(server)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Local Proton Bridge only
	}

	if err := c.StartTLS(tlsConfig); err != nil {
		c.Logout()
		return nil, err
	}

	if err := c.Login(username, password); err != nil {
		c.Logout()
		return nil, err
	}

	return c, nil
}
