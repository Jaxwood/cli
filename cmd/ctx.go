package cmd

import (
	"fmt"

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

		exec.Command("az", "account", "set", "-n", config.SubscriptionId).Start()
		exec.Command("az", "aks", "get-credentials", "-n", config.ClusterName, "-g", config.ResourceGroup, "--overwrite-existing", "--context", env).Start()
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
