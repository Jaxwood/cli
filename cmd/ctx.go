package cmd

import (
	"fmt"
	"bytes"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os/exec"
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
	Long:  `Sets the Azure context using the env flag to specify the environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		var config Azure
		viper.UnmarshalKey(env, &config)
		fmt.Println("Setting environment to: " + env)

		// set the account
		accountCmd := exec.Command("az", "account", "set", "-n", config.SubscriptionId)
		var accountStd, accountErr bytes.Buffer
		accountCmd.Stdout = &accountStd
		accountCmd.Stderr = &accountErr
		err := accountCmd.Run()
		if err != nil {
			   fmt.Println(err)
		    }
		fmt.Println("out:", accountStd.String(), "err:", accountErr.String())

		// update kube config
		kubeCmd := exec.Command("az", "aks", "get-credentials", "-n", config.ClusterName, "-g", config.ResourceGroup, "--overwrite-existing", "--context", env)
		var kubeStd, kubeErr bytes.Buffer
		kubeCmd.Stdout = &kubeStd
		kubeCmd.Stderr = &kubeErr
		err = kubeCmd.Run()
		if err != nil {
			   fmt.Println(err)
		    }
		fmt.Println("out:", kubeStd.String(), "err:", kubeErr.String())
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
