package cmd

import (
  "fmt"
  "bytes"

  "strconv"

  "encoding/json"
  "io/ioutil"
  "net/http"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

func postResource(client *http.Client, token string, url string, image string, tag string) []byte {
  postBody, _ := json.Marshal(map[string]map[string]string{
    "templateParameters": {
      "imageName":  image,
      "imageTag": tag,
      "sourceRepo": "docker-dev",
    },
  })
  responseBody := bytes.NewBuffer(postBody)

  req, err := http.NewRequest(http.MethodPost, url, responseBody)
  req.Header.Add("Authorization", "Basic "+basicAuth("", token))
  req.Header.Add("Content-Type", "application/json")
  response, err := client.Do(req)

  if err != nil {
    fmt.Println(err)
  }
  defer response.Body.Close()
  body, err := ioutil.ReadAll(response.Body)

  if response.StatusCode != http.StatusOK {
    fmt.Println("Non-OK HTTP status:", response.StatusCode)
    panic(err)
  }

  return body
}

var promoteCmd = &cobra.Command{
  Use:   "promote",
  Short: "promote an Docker image",
  Long:  "usage: cli promote --image <sub-folder/image-name> --tag <tag>",
  Run: func(cmd *cobra.Command, args []string) {
    promoteBuildName := viper.GetString("promoteBuildName")

    image, _ := cmd.Flags().GetString("image")
    tag, _ := cmd.Flags().GetString("tag")
    project := viper.GetString("project")
    token := viper.GetString("token")

    client := &http.Client{
      CheckRedirect: redirectPolicyFunc,
    }

    // get the definition by name
    definitionBody := requestResource(client, token, fmt.Sprintf("%s/_apis/build/definitions?name=%s", project, promoteBuildName))
    var definition DefinitionResponse
    json.Unmarshal(definitionBody, &definition)

    // start the promotion build
    id := definition.Value[0].Id
    runBody := postResource(client, token, fmt.Sprintf("%s/_apis/pipelines/%s/runs?api-version=6.0-preview.1", project, strconv.Itoa(id)), image, tag)
    fmt.Println(string(runBody))
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
