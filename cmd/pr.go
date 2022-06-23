/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Repository struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type GitResponse struct {
	Count int          `json:"count"`
	Value []Repository `json:"value"`
}

func getRemote() string {
	gitCmd := exec.Command("git", "remote", "-v")
	var gitStd, gitErr bytes.Buffer
	gitCmd.Stdout = &gitStd
	gitCmd.Stderr = &gitErr
	err := gitCmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	if gitErr.Len() > 0 {
		fmt.Println(gitErr.String())
	}

	return gitStd.String()
}

func getBranchname() string {
	gitCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var gitStd, gitErr bytes.Buffer
	gitCmd.Stdout = &gitStd
	gitCmd.Stderr = &gitErr
	err := gitCmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	if gitErr.Len() > 0 {
		fmt.Println(gitErr.String())
	}

	return strings.Trim(gitStd.String(), "\n")
}

func getRepository(token string, project string) *Repository {
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	gitBody := requestResource(client, token, project+"/_apis/git/repositories?api-version=7.1-preview.1")
	var git GitResponse
	json.Unmarshal(gitBody, &git)

	remotes := strings.Split(getRemote(), "\n")
	r, _ := regexp.Compile("(https://.* )")
	remote := r.FindString(remotes[0])
	segments := strings.Split(remote, "/")
	repositoryName := strings.Trim(segments[len(segments)-1], " ")

	for i := 0; i < len(git.Value); i++ {
		if git.Value[i].Name == repositoryName {
			return &git.Value[i]
		}
	}

	return nil
}

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Open a pull request for the current branch against main",
	Long:  "Support for Azure DevOps",
	Run: func(cmd *cobra.Command, args []string) {
		project := viper.GetString("project")
		token := viper.GetString("token")

		branchName := getBranchname()
		repository := getRepository(token, project)

		url := "https://dev.azure.com/" + project "/_git/" + repository.Name + "/pullrequestcreate?sourceRef=" + branchName + "&targetRef=main&sourceRepositoryId=" + repository.Id + "&targetRepositoryId=" + repository.Id
		exec.Command("open", url).Start()
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
}
