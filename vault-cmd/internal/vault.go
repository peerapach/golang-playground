package internal

import (
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
)

const VAULT_ADDR = "http://localhost:8200"

func NewVault(username, password string) *vault.Client {
	config := vault.DefaultConfig()
	config.Address = VAULT_ADDR

	client, err := vault.NewClient(config)
	if err != nil {
		fmt.Printf("Unable to initialize Vault client: %s \n", err.Error())
		os.Exit(2)
	}

	options := map[string]interface{}{
		"password": password,
	}

	path := fmt.Sprintf("auth/userpass/login/%s", username)

	// PUT call to get a token
	secret, err := client.Logical().Write(path, options)
	if err != nil {
		fmt.Printf("Can't login with Vault, please recheck username or password \n")
		fmt.Printf("Err: %s \n", err.Error())
		os.Exit(2)
	}
	t, err := secret.TokenID()
	if err != nil {
		fmt.Printf("Can't get tokem id after login \n")
		fmt.Printf("%s \n", err.Error())
	}
	client.SetToken(t)

	return client

}

func WriteK8sToken(folderId string, clusterName string, namespace string, token string, client *vault.Client) {
	secretData := map[string]interface{}{
		"data": map[string]interface{}{
			"token": token,
		},
	}

	_, err := client.Logical().Write(folderId+"/data/"+clusterName+"/"+namespace+"/token", secretData)
	if err != nil {
		fmt.Printf("unable to write secret: %v", err.Error())
	}

	fmt.Println("Secret written successfully.")
}

func ReadK8sToken(folderId string, clusterName string, namespace string, client *vault.Client) string {
	v, err := client.Logical().Read(folderId + "/data/" + clusterName + "/" + namespace + "/token")
	if err != nil {
		fmt.Printf("unable to write secret: %v", err)
	}
	t := fmt.Sprintf("%s", v.Data["data"].(map[string]interface{})["token"])
	return t
}
