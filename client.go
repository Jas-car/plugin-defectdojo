package secretsengine

import (
	defectdojo "github.com/truemilk/go-defectdojo/defectdojo"
	"errors"
	"context"
	"encoding/json"
	"ftm"
	"os"
)

// hashiCupsClient creates an object storing
// the client.
type defectdojoClient struct {
	*defectdojo.Client
}

// newClient creates a new client to access HashiCups
// and exposes it for any secrets or roles to use.
func newClient(config *defectdojoConfig, ctx context.Context) (d *defectdojoClient,  err error , myToken string ) {

	if config == nil {
		return nil, errors.New("client configuration was nil"),""
	}

	if config.Username == "" {
		return nil, errors.New("client username was not defined"),""
	}

	if config.Password == "" {
		return nil, errors.New("client password was not defined"),""
	}

	if config.URL == "" {
		return nil, errors.New("client URL was not defined"),""
	}
    
	
	/*dj, _ := defectdojo.NewDojoClient(config.URL, "", nil)

	resp, err := dj.ApiTokenAuth.Create(ctx, &defectdojo.AuthToken{
		Username: &config.Username,
		Password: &config.Password,
	})

	if err != nil {
		return nil, err , ""
	}
	myToken = string(*resp.Token)*/
	url := "https://demo.defectdojo.org"
	//url := "http://host.docker.internal:8080"
	dj, err := defectdojo.NewDojoClient(url, "", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	resp, err := dj.ApiTokenAuth.Create(ctx, &defectdojo.AuthToken{
		Username: defectdojo.Str("admin"),
		Password: defectdojo.Str("1Defectdojo@demo#appsec"),
	})
	if err != nil {
		fmt.Println("main:", err)
		return
	}
	b, err := json.Marshal(resp.Token)
	if err != nil {
		fmt.Println("main:", err)
		return
	}
    fmt.Println(b)
	fmt.Println(string(b))
	return &defectdojoClient{dj}, nil , string(b)
}
