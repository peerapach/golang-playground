/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"ezylinux/vault-cmd/internal"

	"github.com/spf13/cobra"
)

// tokenCmd represents the token command
var tokenReadCmd = &cobra.Command{
	Use:   "token",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Get K8S info
		k8sInfo := ParseYaml(internal.GetFileContentFromGit(K8sInfoRepo, K8sNonProdFileInfo))

		// New k8s client
		clusterHost, _ := internal.NewK8sClient()

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

		t := internal.ReadK8sToken(m, clusterName, namespace, v)
		fmt.Printf("Token: %s \n", t)
	},
}

func init() {
	readCmd.AddCommand(tokenReadCmd)
	tokenReadCmd.Flags().StringVarP(&folderId, "folderid", "f", "", "Secret folder name")
	tokenReadCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace of OpenShift")
	tokenReadCmd.Flags().StringVarP(&vaultUser, "vault_user", "u", "", "A help for foo")

	tokenReadCmd.MarkFlagRequired("folderId")
	tokenReadCmd.MarkFlagRequired("namespace")
	tokenReadCmd.MarkFlagRequired("vaultUser")
}
