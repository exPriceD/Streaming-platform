package clients

import "log"

type Clients struct {
	Auth *AuthClient
}

func NewClients(authServiceAddr string) (*Clients, error) {
	authClient, err := NewAuthClient(authServiceAddr)
	if err != nil {
		log.Printf("Failed to create AuthClient: %v", err)
		return nil, err
	}
	defer func(authClient *AuthClient) {
		err := authClient.Close()
		if err != nil {
			log.Fatalf("Failed to close AuthClient: %v", err)
		}
	}(authClient)

	return &Clients{Auth: authClient}, nil
}
