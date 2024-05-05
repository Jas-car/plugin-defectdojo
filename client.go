package secretsengine

import (
	defectdojo "github.com/truemilk/go-defectdojo/defectdojo"
)

// hashiCupsClient creates an object storing
// the client.
type defectdojoClient struct {
	*defectdojo.Client
}

// newClient creates a new client to access HashiCups
// and exposes it for any secrets or roles to use.
func newClient(config *defectdojoConfig) (*defectdojoClient, error) {

	if config == nil {
		return nil, errors.New("client configuration was nil")
	}

	if config.Username == "" {
		return nil, errors.New("client username was not defined")
	}

	if config.Password == "" {
		return nil, errors.New("client password was not defined")
	}

	if config.URL == "" {
		return nil, errors.New("client URL was not defined")
	}

	dj, _ := defectdojo.NewDojoClient(&config.URL, "", nil)

	resp, err := dj.ApiTokenAuth.Create(ctx, &defectdojo.AuthToken{
		Username: defectdojo.String(&config.Username),
		Password: defectdojo.String(&config.Password),
	})

	if err != nil {
		return nil, err
	}
	return &defectdojoClient{resp}, nil
}
