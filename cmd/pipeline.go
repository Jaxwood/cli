package cmd

import (
  "fmt"
  "os/exec"

  "strconv"

  "encoding/base64"
  "encoding/json"
  "io/ioutil"
  "net/http"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

type DefinitionResponse struct {
  Count int `json:"count"`
  Value []struct {
    Id    int `json:"id"`
    Links struct {
      Web struct {
        Href string `json:"href"`
      } `json:"web"`
    } `json:"_links"`
  } `json:"value"`
}

func basicAuth(username, password string) string {
  auth := username + ":" + password
  return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
  token := viper.GetString("token")
  req.Header.Add("Authorization", "Basic "+basicAuth("", token))
  return nil
}

func requestResource(client *http.Client, token string, url string) []byte {
  req, err := http.NewRequest("GET", url, nil)
  req.Header.Add("Authorization", "Basic "+basicAuth("", token))
  response, err := client.Do(req)

  if err != nil {
    fmt.Println(err)
  }
  defer response.Body.Close()
  body, err := ioutil.ReadAll(response.Body)

  return body
}

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
  Use:   "pipeline",
  Short: "get latest build by name",
  Long:  ``,
  Run: func(cmd *cobra.Command, args []string) {
    name, _ := cmd.Flags().GetString("name")
    project := viper.GetString("project")
    token := viper.GetString("token")

    fmt.Println("Getting build name: " + name)

    client := &http.Client{
      CheckRedirect: redirectPolicyFunc,
    }

    // get the definition by name
    definitionBody := requestResource(client, token, project+"/_apis/build/definitions?name="+name)
    var definition DefinitionResponse
    json.Unmarshal(definitionBody, &definition)

    // get the link to the latest 5 builds
    id := definition.Value[0].Id
    buildBody := requestResource(client, token, project+"/_apis/build/builds?queryOrder=queueTimeDescending&$top=5&definitions="+strconv.Itoa(id))

    var build DefinitionResponse
    json.Unmarshal(buildBody, &build)
    exec.Command("open", build.Value[0].Links.Web.Href).Start()
  },
}

func init() {
  rootCmd.AddCommand(pipelineCmd)

  // Here you will define your flags and configuration settings.
  pipelineCmd.PersistentFlags().StringP("name", "n", "", "Name is required.")
  pipelineCmd.MarkPersistentFlagRequired("name")
}
