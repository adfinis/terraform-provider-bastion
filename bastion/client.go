// Copyright (c) Adfinis
// SPDX-License-Identifier: GPL-3.0-or-later

package bastion

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

var (
	ErrNoAuthMethodsProvided = errors.New("no authentication method provided")
	ErrMissingConfig         = errors.New("missing configuration")
	ErrHostRequired          = errors.New("host is required")
	ErrInvalidPort           = errors.New("invalid port")
	ErrUsernameRequired      = errors.New("username is required")
)

type Config struct {
	Host                  string
	Port                  int
	Username              string
	Timeout               int
	StrictHostKeyChecking bool
}

type Client struct {
	Host         string
	Port         int
	sshClientCfg *ssh.ClientConfig
}

func New(cfg *Config, authMethods ...SSHAuthMethod) (*Client, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	if len(authMethods) == 0 {
		return nil, ErrNoAuthMethodsProvided
	}

	var methods []ssh.AuthMethod
	for _, auth := range authMethods {
		method, err := auth()
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	sshCfg := &ssh.ClientConfig{
		User:            cfg.Username,
		Auth:            methods,
		HostKeyCallback: getHostKeyCallback(cfg.StrictHostKeyChecking),
	}

	return &Client{
		Host:         cfg.Host,
		Port:         cfg.Port,
		sshClientCfg: sshCfg,
	}, nil
}

// validateConfig checks that the provided configuration is valid.
func validateConfig(cfg *Config) error {
	if cfg == nil {
		return ErrMissingConfig
	}
	if cfg.Host == "" {
		return ErrHostRequired
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return ErrInvalidPort
	}

	if cfg.Username == "" {
		return ErrUsernameRequired
	}

	return nil
}

// getHostKeyCallback returns appropriate host key callback.
func getHostKeyCallback(strictHostKeyChecking bool) ssh.HostKeyCallback {
	if strictHostKeyChecking {
		return nil
	}
	return ssh.InsecureIgnoreHostKey()
}

// sshClient returns a new ssh.Client based on the provided configuration.
func (c *Client) sshClient() (*ssh.Client, error) {
	if c.sshClientCfg == nil {
		return nil, ErrMissingConfig
	}
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return ssh.Dial("tcp", address, c.sshClientCfg)
}
