package cmd

import (
  "bytes"
  "fmt"
  "os"
  "os/exec"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

type Cosign struct {
  ClientId string `mapstructure:"clientid"`
  ClientSecret   string `mapstructure:"clientsecret"`
  Keyvault string `mapstructure:"keyvault"`
  Keyname string `mapstructure:"keyname"`
}

// cosignCmd represents the sysdig command
var cosignCmd = &cobra.Command{
  Use:   "cosign",
  Short: "sign an image using cosign cli",
  Long:  "usage: cli cosign --image org.jfrog.io/docker/sub-folder-image-name:tag",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("cosign called")
    image, _ := cmd.Flags().GetString("image")

    var cosignConfig Cosign
    viper.UnmarshalKey("cosign", &cosignConfig)
    tenant := viper.GetString("tenant")

    // sign an image
    os.Setenv("KEY_VAULT", cosignConfig.Keyvault)
    os.Setenv("AZURE_CLIENT_ID", cosignConfig.ClientId)
    os.Setenv("AZURE_CLIENT_SECRET", cosignConfig.ClientSecret)
    os.Setenv("AZURE_TENANT_ID", tenant)
    signCmd := exec.Command("cosign", "sign", "--key", "azurekms://" + cosignConfig.Keyvault + ".vault.azure.net/" + cosignConfig.Keyname, image)

    var signStd, signStdErr bytes.Buffer
    signCmd.Stdout = &signStd
    signCmd.Stderr = &signStdErr
    err := signCmd.Run()
    if err != nil {
      fmt.Println(err)
    }
    fmt.Println("out:", signStd.String(), "err:", signStdErr.String())
  },
}

func init() {
  rootCmd.AddCommand(cosignCmd)
  cosignCmd.PersistentFlags().StringP("image", "i", "", "Image to sign is required")
  cosignCmd.MarkPersistentFlagRequired("image")
}
