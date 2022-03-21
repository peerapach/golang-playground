/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"ezylinux/vault-cmd/internal"

	promptui "github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// tokenCmd represents the token command
var tokenAddCmd = &cobra.Command{
	Use:   "token",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Get K8S info
		k8sInfo := ParseYaml(internal.GetFileContentFromGit(K8sInfoRepo, K8sNonProdFileInfo))

		// New k8s client
		clusterHost, c := internal.NewK8sClient()
		// Get token by service account
		token := internal.GetTokenServiceAccount(serviceAccount, namespace, c)

		//Get cluster name
		clusterName := getClusterName(clusterHost, k8sInfo)

		// Input password for login to vault
		password := inputPassword()

		// New Vault client
		v := internal.NewVault(vaultUser, password)
		// check KV
		ls, err := v.Sys().ListMounts()
		if err != nil {
			fmt.Printf("Can't list KV from Vault \n")
			fmt.Printf("%s", err.Error())
		}
		m := selectAppFolder(ls)

		if len(internal.ReadK8sToken(m, clusterName, namespace, v)) > 0 {
			prompt := promptui.Prompt{
				Label:     "Token already exist, Over-write Token Confirm",
				IsConfirm: true,
			}
			result, _ := prompt.Run()

			if strings.ToLower(result) == "y" {
				internal.WriteK8sToken(m, clusterName, namespace, token, v)
			}

		} else {
			internal.WriteK8sToken(m, clusterName, namespace, token, v)
		}
	},
}

func init() {
	addCmd.AddCommand(tokenAddCmd)
	tokenAddCmd.Flags().StringVarP(&folderId, "folderid", "f", "", "Secret folder name")
	tokenAddCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace of OpenShift")
	tokenAddCmd.Flags().StringVarP(&serviceAccount, "serviceaccount", "s", "", "Serviceaccount name for ")
	tokenAddCmd.Flags().StringVarP(&vaultUser, "vault_user", "u", "", "Vault user")

	tokenAddCmd.MarkFlagRequired("folderId")
	tokenAddCmd.MarkFlagRequired("namespace")
	tokenAddCmd.MarkFlagRequired("serviceaccount")
	tokenAddCmd.MarkFlagRequired("vaultUser")
}
