package secretsengine

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Factory returns a new backend as logical.Backend
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// hashiCupsBackend defines an object that
// extends the Vault backend and stores the
// target API's client.
type defectdojoBackend struct {
	*framework.Backend
	lock   sync.RWMutex
	client *defectdojoClient
}

// backend defines the target API backend
// for Vault. It must include each path
// and the secrets it will store.
func backend() *defectdojoBackend {
	var b = defectdojoBackend{}

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{},
			SealWrapStorage: []string{
				"config",
				"role/*",
			},
		},
		Paths: framework.PathAppend(
			pathRole(&b),
			[]*framework.Path{
				pathConfig(&b),
				pathCredentials(&b),
			},
		),
		Secrets: []*framework.Secret{
			b.defectdojoToken(),
		},
		BackendType: logical.TypeLogical,
		Invalidate:  b.invalidate,
	}
	return &b
}

// reset clears any client configuration for a new
// backend to be configured
func (b *defectdojoBackend) reset() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.client = nil
}

// invalidate clears an existing client configuration in
// the backend
func (b *defectdojoBackend) invalidate(ctx context.Context, key string) {
	if key == "config" {
		b.reset()
	}
}

// getClient locks the backend as it configures and creates a
// a new client for the target API
func (b *defectdojoBackend) getClient(ctx context.Context, s logical.Storage) (d *defectdojoClient, err error, myToken string ) {
	b.lock.RLock()
	unlockFunc := b.lock.RUnlock
	defer func() { unlockFunc() }()

	if b.client != nil {
		return b.client, nil,""
	}

	b.lock.RUnlock()
	b.lock.Lock()
	unlockFunc = b.lock.Unlock

	config, err := getConfig(ctx, s)
	if err != nil {
		return nil, err,""
	}

	if config == nil {
		config = new(defectdojoConfig)
	}

	b.client, err , myToken = newClient(config, ctx)
	if err != nil {
		return nil, err,""
	}

	return b.client, nil,myToken
}

// backendHelp should contain help information for the backend
const backendHelp = `
The Defectdojo secrets backend dynamically generates user tokens.
After mounting this backend, credentials to manage HashiCups user tokens
must be configured with the "config/" endpoints.
`
