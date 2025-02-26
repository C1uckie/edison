package main

import (
	"code.gitea.io/sdk/gitea"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const version_constraint = "1.23.3"

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

func main() {
	var user_configuration = load_config()
	var client = create_gitea_client(user_configuration.URI, user_configuration.Token)

	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			print_help()
			return
		}
		if arg == "-g" || arg == "--gitea" {
			print_gitea_server_version(client)
			return
		}
		if arg == "-c" || arg == "--config" {
			print_config(user_configuration)
			return
		}
		if arg == "-u" || arg == "--user" {
			print_gitea_user(client)
			return
		}
		if arg == "-r" || arg == "--repos" {
			print_user_repos(client)
			return
		}
		if arg == "-o" || arg == "--orgs" {
			print_org_repos(client)
			return
		}
	}

	print_edison_fetch()
}
func print_edison_fetch() {
	fmt.Println("Default Fetch, should listen to config.json")
}

func print_help() {
	fmt.Println("-h or --help for this menu")
	fmt.Println("-g or --gitea for gitea version")
	fmt.Println("-c or --config to see config")
	fmt.Println("-u or --user to see gitea user")
	fmt.Println("-r or --repos to see user repos")
	fmt.Println("-o or --orgs to see org repos")
}

func load_config() Configuration {
	jsonFile, err := os.Open("./configuration/configuration.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var user_configuration Configuration
	json.Unmarshal(byteValue, &user_configuration)
	return user_configuration
}

func create_gitea_client(config_uri string, config_token string) *gitea.Client {
	client, err := gitea.NewClient(config_uri, gitea.SetToken(config_token))
	if err != nil {
		fmt.Println("Error creating Gitea client:", err)
	}
	return client
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

func print_gitea_server_version(client *gitea.Client) {
	version, _, err := client.ServerVersion()
	if err != nil {
		fmt.Println("Error getting Gitea server version:", err)
	}
	fmt.Println(version)
}

func print_gitea_user(client *gitea.Client) {
	gitea_user_data, _, err := client.GetMyUserInfo()
	if err != nil {
		fmt.Println("Error getting Gitea user:", err)
	}

	fmt.Println("Username:", gitea_user_data.UserName)
}

func print_user_repos(client *gitea.Client) {
	user_repos, _, err := client.ListMyRepos(gitea.ListReposOptions{})
	if err != nil {
		fmt.Println("Error getting Gitea user repos:", err)
	}
	for _, repo := range user_repos {
		gitea_user_data, _, err := client.GetMyUserInfo()
		if err != nil {
			fmt.Println("Error getting Gitea user:", err)
		}
		if repo.Owner.UserName == gitea_user_data.UserName {
			fmt.Println(repo.Name)
			fmt.Println(repo.SSHURL)
		}
	}
}
func print_org_repos(client *gitea.Client) {
	user_orgs, _, err := client.ListMyOrgs(gitea.ListOrgsOptions{})
	if err != nil {
		fmt.Println("Error getting Gitea user orgs:", err)
	}
	for _, org := range user_orgs {
		var org_name = org.UserName
		fmt.Println(org.UserName)
		org_repos, _, err := client.ListOrgRepos(org_name, gitea.ListOrgReposOptions{})
		if err != nil {
			fmt.Println("Error getting Gitea org repos:", err)
		}
		for _, repo := range org_repos {
			fmt.Println(repo.Name)
			fmt.Println(repo.SSHURL)
		}
	}

}
