package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"strings"
)

type Token struct {
  AccessToken  string `json:"access_token"`
  ExpiresIn    string `json:"expires_in"`
  ExtExpiresIn string `json:"ext_expires_in"`
  TokenType    string `json:"token_type"`
}

type Application struct {
  Dev OAuth `mapstructure:"dev"`
  Tst OAuth `mapstructure:"tst"`
  Prd OAuth `mapstructure:"prd"`
}

type OAuth struct {
  Scope        string `mapstructure:"scope"`
  ClientId     string `mapstructure:"client_id"`
  ClientSecret string `mapstructure:"client_secret"`
  GrantType    string `mapstructure:"grant_type"`
}

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
  Use:   "token",
  Short: "Get an access token for application",
  Long: "usage: cli token -e <env> -a <app>",
  Run: func(cmd *cobra.Command, args []string) {
    // read config
    env, _ := cmd.Flags().GetString("env")
    app, _ := cmd.Flags().GetString("app")
    var config Application
    if app == "customer" {
      viper.UnmarshalKey("customersegmentation", &config)
    }
    if app == "asset" {
      viper.UnmarshalKey("assetsync", &config)
    }
    tenant := viper.GetString("tenant")

    var oauth OAuth
    if env == "prd" {
      oauth = config.Prd
    } else {
      oauth = config.Dev
    }
    // request token
    data := url.Values{}
    data.Set("client_id", oauth.ClientId)
    data.Set("client_secret", oauth.ClientSecret)
    data.Set("scope", oauth.Scope)
    data.Set("grant_type", oauth.GrantType)
    response, err := http.Post("https://login.microsoftonline.com/" + tenant + "/oauth2/v2.0/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
    if err != nil {
      fmt.Println(err)
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)

    if err != nil {
      fmt.Println(err)
    }

    var token Token
    json.Unmarshal(body, &token) 
    fmt.Println(token.AccessToken)
  },
}

func init() {
  rootCmd.AddCommand(tokenCmd)

  tokenCmd.PersistentFlags().StringP("env", "e", "", "Env is required. Supported values dev|tst|prd")
  tokenCmd.MarkPersistentFlagRequired("env")

  tokenCmd.PersistentFlags().StringP("app", "a", "", "Application is required. Supported values customer|asset")
  tokenCmd.MarkPersistentFlagRequired("app")

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // tokenCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // tokenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
