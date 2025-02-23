package clients

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Clients struct {
	Auth *AuthClient
}

type Config struct {
	AuthServiceAddr string
	DialTimeout     time.Duration
}

func NewClients(cfg Config) (*Clients, error) {
	authCfg := AuthClientConfig{
		Address:     cfg.AuthServiceAddr,
		DialTimeout: cfg.DialTimeout,
		UseTLS:      false,
	}
	authClient, err := NewAuthClient(authCfg)
	if err != nil {
		return nil, fmt.Errorf("initialize auth client: %w", err)
	}
	return &Clients{
		Auth: authClient,
	}, nil
}

func (c *Clients) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.Auth.Close(); err != nil {
			errChan <- fmt.Errorf("close auth client: %w", err)
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
