package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bytes"
	"fmt"
	"os/exec"
)

var promoteCmd = &cobra.Command{
  Use:   "promote",
  Short: "promote an Docker image",
  Long:  "usage: cli promote --image <sub-folder/image-name> --tag <tag>",
  Run: func(cmd *cobra.Command, args []string) {
    promoteBuildName := viper.GetString("promoteBuildName")
    promoteProject := viper.GetString("promoteProject")
    image, _ := cmd.Flags().GetString("image")
    tag, _ := cmd.Flags().GetString("tag")
    promoteExec := exec.Command("gh", "workflow", "run", "-R", promoteProject, "-f", "image_name="+image, "-f", "image_tag="+tag, promoteBuildName)
    var stdOut, stdErr bytes.Buffer
    promoteExec.Stdout = &stdOut
    promoteExec.Stderr = &stdErr
    if err := promoteExec.Run(); err != nil {
      fmt.Println(err)
    }
  },
}

func init() {
  rootCmd.AddCommand(promoteCmd)

  // Here you will define your flags and configuration settings.
  promoteCmd.PersistentFlags().StringP("image", "i", "", "Image is required.")
  promoteCmd.MarkPersistentFlagRequired("image")

  promoteCmd.PersistentFlags().StringP("tag", "t", "", "Tag is required.")
  promoteCmd.MarkPersistentFlagRequired("tag")
}
