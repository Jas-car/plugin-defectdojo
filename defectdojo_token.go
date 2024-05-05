package secretsengine

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	defectdojoTokenType = "defectdojo_token"
)

// hashiCupsToken defines a secret for the HashiCups token
type defectdojoToken struct {
	Username string `json:"username"`
	TokenID  string `json:"token_id"`
	Token    string `json:"token"`
}


// hashiCupsToken defines a secret to store for a given role
// and how it should be revoked or renewed.
func (b *defectdojoBackend) defectdojoToken() *framework.Secret {
	return &framework.Secret{
		Type: defectdojoTokenType,
		Fields: map[string]*framework.FieldSchema{
			"token": {
				Type:        framework.TypeString,
				Description: "HashiCups Token",
			},
		},
	}
}

// createToken calls the HashiCups client to sign in and returns a new token
func createToken(ctx context.Context, c *defectdojoClient, username string) (*defectdojoToken, error) {
	response, err := *c.Token
	//resp, err := c.ApiTokenAuth.Create(ctx, &defectdojo.AuthToken{
	//	Username: defectdojo.String(&config.Username),
	//	Password: defectdojo.String(&config.Password),
	//})	
	if err != nil {
		return nil, fmt.Errorf("error creating HashiCups token: %w", err)
	}

	tokenID := uuid.New().String()

	return &defectdojoToken{
		Username: username,
		TokenID:  tokenID,
		Token:    *c.Token,
	}, nil
}