package main

import (
	"code.gitea.io/sdk/gitea"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Configuration struct {
	Ascii_Art      string `json:"ascii_art"`
	Top_Langs      int    `json:"top_langs"`
	Repo_Count     bool   `json:"repo_count"`
	Include_Orgs   bool   `json:"include_orgs"`
	Gitea_User     bool   `json:"gitea_user"`
	Gitea_Version  bool   `json:"gitea_version"`
	Edison_Version bool   `json:"edison_version"`
	Token          string `json:"token"`
	URI            string `json:"URI"`
}

func load_config() Configuration {
	jsonFile, err := os.Open("./configuration/configuration.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened ./configuration/configuration.json")
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var user_configuration Configuration
	json.Unmarshal(byteValue, &user_configuration)
	return user_configuration
}

func print_config(user_configuration Configuration) {
	fmt.Println("Token: ", user_configuration.Token)
	fmt.Println("URI: ", user_configuration.URI)
	fmt.Println("Ascii Art: ", user_configuration.Ascii_Art)
	fmt.Println("Top Langs: ", user_configuration.Top_Langs)
	fmt.Println("Repo Count: ", user_configuration.Repo_Count)
	fmt.Println("Include Orgs: ", user_configuration.Include_Orgs)
	fmt.Println("Gitea User: ", user_configuration.Gitea_User)
	fmt.Println("Gitea Version: ", user_configuration.Gitea_Version)
	fmt.Println("Edison Version: ", user_configuration.Edison_Version)
}

// WIP
/*
func create_gitea_client() {
	client, err := gitea.NewClient("http://192.168.7.2:3029", gitea.SetToken("4134b116be73b40c8bc8051dd29fc76f64d53f23"))
	if err != nil {
		fmt.Println("Error creating Gitea client:", err)
		return
	}
}
*/

func main() {
	var user_configuration = load_config()
	print_config(user_configuration)
}
