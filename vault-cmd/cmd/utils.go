package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	promptui "github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

var folderId, namespace, serviceAccount, vaultUser string

type ClusterInfo struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Ip      string `yaml:"ip"`
	Version string `yaml:"version"`
	Console string `yaml:"console"`
}

func inputPassword() string {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("password must have more than 6 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Password",
		Validate: validate,
		Mask:     '*',
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return result
}

func selectAppFolder(o map[string]*api.MountOutput) string {
	var m []string
	for i := range o {
		if strings.Contains(i, folderId) {
			s := strings.TrimSuffix(i, "/")
			m = append(m, s)
		}
	}

	if len(m) == 0 {
		fmt.Printf("Can't find App folder from Vault")
		os.Exit(2)

	} else if len(m) > 1 {
		prompt := promptui.Select{
			Label: "Please Select Folder ID",
			Items: m,
		}
		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(2)
		}
		return result
	}

	return m[0]
}

func getClusterName(clusterHost string, k8sInfo map[string]ClusterInfo) string {
	for k, v := range k8sInfo {
		if clusterHost == "https://"+v.Host+":"+v.Port {
			clusterName := k
			return clusterName
		}
	}
	return ""
}

func ParseYaml(f []byte) map[string]ClusterInfo {

	m := make(map[string]ClusterInfo)

	err := yaml.Unmarshal(f, &m)
	if err != nil {
		fmt.Printf("Can't parse yaml: %v", err.Error())
		os.Exit(2)
	}

	return m
}
