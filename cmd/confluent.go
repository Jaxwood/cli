package cmd

import (
  "bytes"
  "fmt"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"

  "os/exec"
)

type Confluent struct {
  Dev State `mapstructure:"dev"`
  Tst State `mapstructure:"tst"`
  Reg State `mapstructure:"reg"`
  Prd State `mapstructure:"prd"`
}

type State struct {
  ResourceGroupName  string `mapstructure:"resource_group_name"`
  StorageAccountName string `mapstructure:"storage_account_name"`
  ContainerName      string `mapstructure:"container_name"`
  Key                string `mapstructure:"key"`
  AccessKey          string `mapstructure:"access_key"`
}

// confluentCmd represents the confluent command
var confluentCmd = &cobra.Command{
  Use:   "confluent",
  Short: "use to run terraform init for a given environment",
  Long: `supported environments: dev, tst and prd`,
  Run: func(cmd *cobra.Command, args []string) {
    env, _ := cmd.Flags().GetString("env")
    var config Confluent
    viper.UnmarshalKey("confluent", &config)

    var state State
    if env == "prd" {
      state = config.Prd
    } else if env == "tst" {
      state = config.Tst
    } else if env == "reg" {
      state = config.Reg
    } else {
      state = config.Dev
    }
    resource_group_arg := "-backend-config=resource_group_name=" + state.ResourceGroupName
    storage_account_arg := "-backend-config=storage_account_name=" + state.StorageAccountName
    container_name_arg := "-backend-config=container_name=" + state.ContainerName
    key_arg := "-backend-config=key=" + state.Key
    access_key_arg := "-backend-config=access_key=" + state.AccessKey
    outCmd := exec.Command("terraform", "init", resource_group_arg, storage_account_arg, container_name_arg, key_arg, access_key_arg)
    var outb, errb bytes.Buffer
    outCmd.Stdout = &outb
    outCmd.Stderr = &errb
    err := outCmd.Run()
    if err != nil {
         fmt.Println(err)
        }
    fmt.Println("out:", outb.String(), "err:", errb.String())
  },
}

func init() {
  rootCmd.AddCommand(confluentCmd)
  confluentCmd.PersistentFlags().StringP("env", "e", "", "Env is required.")
  confluentCmd.MarkPersistentFlagRequired("env")
}
