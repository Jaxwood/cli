package cmd

import (
	"fmt"
	"bytes"

	"os/exec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "get latest build by name",
	Long: `requires the az devops extentions to be installed for azure cli`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		project := viper.GetString("project")
		fmt.Println("Getting build name: " + name)

		// get the definition id
		definitionCmd := exec.Command("az", "pipelines", "list", "--query", "[?name=='" + name + "'].id | [0]")
		var definitionOut, definitionErr bytes.Buffer
		definitionCmd.Stdout = &definitionOut
		definitionCmd.Stderr = &definitionErr
		_ = definitionCmd.Run()

		// get the last run for that build
		buildCmd := exec.Command("az", "pipelines", "build", "list", "--query", "sort_by([?definition.id == `" + definitionOut.String() + "`], &queueTime)[-1].id")
		var outb, errb bytes.Buffer
		buildCmd.Stdout = &outb
		buildCmd.Stderr = &errb
		_ = buildCmd.Run()
		exec.Command("open", project + "/_build/results?view=results&buildId=" + outb.String()).Start()
	},
}

func init() {
	rootCmd.AddCommand(pipelineCmd)

	// Here you will define your flags and configuration settings.
	pipelineCmd.PersistentFlags().StringP("name", "n", "", "Name is required.")
	pipelineCmd.MarkPersistentFlagRequired("name")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipelineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pipelineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
