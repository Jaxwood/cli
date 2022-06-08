package cmd

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Sysdig struct {
	Token string `mapstructure:"token"`
	Url   string `mapstructure:"url"`
	Image string `mapstructure:"image"`
}

type Jfrog struct {
	Username string `mapstructure:"username"`
	Token    string `mapstructure:"token"`
}

// sysdigCmd represents the sysdig command
var sysdigCmd = &cobra.Command{
	Use:   "sysdig",
	Short: "scan an image using sysdig inline scanner",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sysdig called")
		image, _ := cmd.Flags().GetString("image")

		var sysdigConfig Sysdig
		viper.UnmarshalKey("sysdig", &sysdigConfig)

		var jfrogConfig Jfrog
		viper.UnmarshalKey("jfrog", &jfrogConfig)

		// scan image
		dockerCmd := exec.Command("nerdctl", "run", sysdigConfig.Image, "-k", sysdigConfig.Token, "-s", sysdigConfig.Url, "--registry-auth-basic", jfrogConfig.Username+":"+jfrogConfig.Token, image)

		var dockerStd, dockerStdErr bytes.Buffer
		dockerCmd.Stdout = &dockerStd
		dockerCmd.Stderr = &dockerStdErr
		err := dockerCmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("out:", dockerStd.String(), "err:", dockerStdErr.String())
	},
}

func init() {
	rootCmd.AddCommand(sysdigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sysdigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	sysdigCmd.PersistentFlags().StringP("image", "i", "", "Image to scan is required")
	sysdigCmd.MarkPersistentFlagRequired("image")
}
