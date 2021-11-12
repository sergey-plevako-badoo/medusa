package vaultengine

import (
	"github.com/jonasvinther/medusa/pkg/importer"
	"testing"
)

const (
	address = "https://127.0.0.1:8201"
	token   = "00000000-0000-0000-0000-000000000000"
)

var client *Client

func init() {
	client = NewClient(address, token, true, "sys")
}

func TestCreateKvEngine(t *testing.T) {
	err := client.EnableSecretsEngine("k8s")
	if err != nil {
		t.Error("Unable to add Secrets kv engine")
	}
}

func TestEnableAuth(t *testing.T) {
	errU := client.EnableUserPass()
	if errU != nil {
		t.Error("Cannot enable userpass Auth method")
	}

	errK := client.EnableKubernetes()
	if errK != nil {
		t.Error("Cannot enable kuber Auth method")
	}
}

func TestAddUser(t *testing.T) {
	_, err := client.AddUser("test", "tester", []string{"default"})
	if err != nil {
		t.Error("Unable to add a user")
	}
}

var usersYaml = []byte(`entity:
  name:
    ent-dev-tester:
      name: ent-dev-tester
      policies:
      - default
`)

func TestImportUsers(t *testing.T) {
	parsedYaml, errParse := importer.Import(usersYaml)

	if errParse != nil {
		t.Error("Failed parsing yaml")
	}

	err := client.ImportUsers(parsedYaml)
	if err != nil {
		t.Error("Failed importing user from yaml")
	}
}
