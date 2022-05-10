package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os/exec"
)

// ctxCmd represents the ctx command
var ctxCmd = &cobra.Command{
	Use:   "ctx",
	Short: "CCoE",
	Long:  `CCoE developer productivity tool`,
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		fmt.Println("Setting environment to: " + env)
		account_id := viper.GetString(env)
		fmt.Println("Setting environment to: " + account_id)
		exec.Command("az", "account", "set", "-n", account_id).Start()
	},
}

func init() {
	rootCmd.AddCommand(ctxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	ctxCmd.PersistentFlags().StringP("env", "e", "", "Env is required. Supported values dev|tst|prd")
	ctxCmd.MarkPersistentFlagRequired("env")
}
