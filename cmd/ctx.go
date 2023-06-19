package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os/exec"

	"github.com/jaxwood/cli/internal/types"
)

type Azure struct {
  ResourceGroup  string `mapstructure:"resource_group"`
  ClusterName    string `mapstructure:"cluster_name"`
  SubscriptionId string `mapstructure:"subscription_id"`
}

// ctxCmd represents the ctx command
var ctxCmd = &cobra.Command{
  Use:   "ctx",
  Short: "Sets the Azure context",
  Long:  "usage: cli ctx -e <env>",
  Run: func(cmd *cobra.Command, args []string) {
    env, _ := cmd.Flags().GetString("env")
    var config types.Azure
    viper.UnmarshalKey(env, &config)
    fmt.Println("Setting environment to " + env + "..")

    // set the account
    accountCmd := exec.Command("az", "account", "set", "-n", config.SubscriptionId)
    var accountStd, accountErr bytes.Buffer
    accountCmd.Stdout = &accountStd
    accountCmd.Stderr = &accountErr
    err := accountCmd.Run()
    if err != nil {
         fmt.Println(err)
        }
    if accountStd.Len() > 0 || accountErr.Len() > 0 {
      fmt.Println("out:", accountStd.String(), "err:",  accountErr.String())
    }

    // update kube config
    fmt.Println("Updating kubeconfig..")
    kubeCmd := exec.Command("az", "aks", "get-credentials", "-n", config.ClusterName, "-g", config.ResourceGroup, "--overwrite-existing", "--context", env)
    var kubeStd, kubeErr bytes.Buffer
    kubeCmd.Stdout = &kubeStd
    kubeCmd.Stderr = &kubeErr
    err = kubeCmd.Run()
    if err != nil {
         fmt.Println(err)
        }
    if kubeStd.Len() > 0 || kubeErr.Len() > 0 {
      fmt.Println("out:", kubeStd.String(), "err:", kubeErr.String())
    }

    // remove kubelogin tokens
    fmt.Println("Removing kubelogin tokens..")
    kubetokenCmd := exec.Command("kubelogin", "remove-tokens")
    var kubetokenStd, kubetokenErr bytes.Buffer
    kubetokenCmd.Stdout = &kubetokenStd
    kubetokenCmd.Stderr = &kubetokenErr
    err = kubetokenCmd.Run()
    if err != nil {
         fmt.Println(err)
        }
    if kubetokenStd.Len() > 0 || kubetokenErr.Len() > 0 {
      fmt.Println("out:", kubetokenStd.String(), "err:", kubetokenErr.String())
    }

    // run kubelogin
    fmt.Println("Running kubelogin..")
    kubeloginCmd := exec.Command("kubelogin", "convert-kubeconfig", "-l", "azurecli")
    var kubeloginStd, kubeloginErr bytes.Buffer
    kubeloginCmd.Stdout = &kubeloginStd
    kubeloginCmd.Stderr = &kubeloginErr
    err = kubeloginCmd.Run()
    if err != nil {
         fmt.Println(err)
        }
    if kubeloginStd.Len() > 0 || kubeloginErr.Len() > 0 {
      fmt.Println("out:", kubeloginStd.String(), "err:", kubeloginErr.String())
    }
  },
}

func init() {
  rootCmd.AddCommand(ctxCmd)

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  ctxCmd.PersistentFlags().StringP("env", "e", "", "Env is required. Supported values (dev | tst | prd)")
  ctxCmd.MarkPersistentFlagRequired("env")
}
