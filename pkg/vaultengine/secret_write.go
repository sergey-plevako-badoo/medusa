package vaultengine

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
)

// SecretWrite is used for writing data to a Vault instance
func (client *Client) SecretWrite(path string, data map[string]interface{}) {
	infix := "/data/"

	if client.engineType == "kv1" {
		infix = "/"
	}

	finalPath := client.engine + infix + path

	finalData := make(map[string]interface{})

	if client.engineType == "kv1" {
		finalData = data
	} else {
		finalData["data"] = data
	}

	// If the data object has the json-object key
	// it means that the secret is not in the default
	// key/value format.
	if jsonVal, ok := data["json-object"]; ok {
		var jsonString string

		// The kv2 engine needs the data wrapped in a "data" key
		if client.engineType == "kv2" {
			jsonString = fmt.Sprintf("{\"data\":%s}", jsonVal)
		} else {
			jsonString = jsonVal.(string)
		}

		_, err := client.vc.Logical().WriteBytes(finalPath, []byte(jsonString))
		if err != nil {
			fmt.Printf("Error while writing secret. %s\n", err)
		} else {
			fmt.Printf("Secret successfully written to Vault [%s] using path [%s]\n", client.addr, path)
		}
	} else {
		_, err := client.vc.Logical().Write(finalPath, finalData)
		if err != nil {
			fmt.Printf("Error while writing secret. %s\n", err)
		} else {
			fmt.Printf("Secret successfully written to Vault [%s] using path [%s]\n", client.addr, path)
		}
	}
}

// Create k8s secrets mount v2
func (client *Client) EnableSecretsEngine(path string) error {
	options := map[string]string{"version": "2"}
	input := &vault.MountInput{
		Type:    "kv",
		Options: options,
	}

	return client.vc.Sys().Mount(path, input)
}

func (client *Client) EnableUserPass() error {
	options := &vault.EnableAuthOptions{
		Type: "userpass",
	}
	return client.vc.Sys().EnableAuthWithOptions("userpass", options)
}

func (client *Client) AddUser(name string, pass string, policies []string) (string, error) {
	params := map[string]interface{}{
		"name":     name,
		"password": pass,
		"policies": policies, // the name of the role in Vault that was created with this app's Kubernetes service account bound to it
	}

	r := client.vc.NewRequest("POST", "/v1/auth/userpass/users/"+name)
	if err := r.SetJSONBody(params); err != nil {
		return "", err
	}

	resp, err := client.vc.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// any 404 indicates k/v v1
		if resp != nil && resp.StatusCode == 404 {
			return "", nil
		}
		return "", err
	}

	return "", err
}

func (client *Client) AddUserAlias() error {
	return nil
}
