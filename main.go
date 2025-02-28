package main

import (
	"bufio"
	"code.gitea.io/sdk/gitea"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const version_constraint = "1.23.3"

type Configuration struct {
	Ascii_Art      []string `json:"ascii_art"`
	Top_Langs      int      `json:"top_langs"`
	Repo_Count     bool     `json:"repo_count"`
	Include_Orgs   bool     `json:"include_orgs"`
	Gitea_User     bool     `json:"gitea_user"`
	Gitea_Version  bool     `json:"gitea_version"`
	Edison_Version bool     `json:"edison_version"`
	Token          string   `json:"token"`
	URI            string   `json:"URI"`
}

type langStat struct {
	Lang       string
	Percentage float64
}

type Color string

const (
	AutumnGreen  Color = "AutumnGreen"
	AutumnRed    Color = "AutumnRed"
	AutumnYellow Color = "AutumnYellow"
	DefaultReset Color = "DefaultReset"
)

// Some Kanagawa.nvim colors
var ColorMap = map[Color]string{
	AutumnGreen:  "\x1B[38;2;118;148;106m",
	AutumnRed:    "\x1B[38;2;195;64;67m",
	AutumnYellow: "\x1B[38;2;220;165;97m",
	DefaultReset: "\x1B[0m",
}

func main() {
	make_config_dir()
	var user_configuration = load_config()
	var client = create_gitea_client(user_configuration.URI, user_configuration.Token)
	args := os.Args[1:]
	recognizedFlag := false

	if len(args) == 0 {
		// default fetch
		print_ascii_art(user_configuration.Ascii_Art)
		print_gitea_server_version(client)
		print_gitea_user(client)
		print_user_repo_count(user_configuration.Repo_Count, user_configuration.Include_Orgs, client)
		print_user_total_loc(client, user_configuration.Include_Orgs)
		print_user_langs(client, user_configuration.Include_Orgs, true, user_configuration.Top_Langs)
		return
	}

	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			print_help()
			return
		}
		if arg == "-v" || arg == "--version" {
			print_edison_version()
			return
		}
		if arg == "-pa" || arg == "--print-ascii" {
			print_ascii_art(user_configuration.Ascii_Art)
			return
		}
		if arg == "-gv" || arg == "--gitea-version" {
			print_gitea_server_version(client)
			return
		}
		if arg == "-pc" || arg == "--print-config" {
			print_config(user_configuration)
			return
		}
		if arg == "-gu" || arg == "--get-user" {
			print_gitea_user(client)
			return
		}
		if arg == "-lr" || arg == "--list-repos" {
			print_user_repos(client)
			return
		}
		if arg == "-lor" || arg == "--list-org-repos" {
			print_org_repos(client)
			return
		}
		if arg == "-rc" || arg == "--repo-count" {
			print_user_repo_count(user_configuration.Repo_Count, user_configuration.Include_Orgs, client)
			return
		}
		if arg == "-cr" || arg == "--create-repo" {
			create_gitea_repo(client)
			return
		}
		if arg == "-pul" || arg == "--print-user-langs" {
			print_user_langs(client, user_configuration.Include_Orgs, false, user_configuration.Top_Langs)
		}
		if arg == "-tul" || arg == "--top-user-langs" {
			print_user_langs(client, user_configuration.Include_Orgs, true, user_configuration.Top_Langs)
		}
		if arg == "-loc" || arg == "--lines-of-code" {
			print_user_total_loc(client, user_configuration.Include_Orgs)
		}
		recognizedFlag = true
	}
	if !recognizedFlag {
		fmt.Println(ColorMap[AutumnRed] + "Error: Unrecognized flag." + ColorMap[DefaultReset])
		print_help()
		return
	}
}

func print_help() {
	fmt.Println(ColorMap[AutumnGreen] + "Usage:" + ColorMap[AutumnYellow] + "      edison [flags]" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "No Flags:" + ColorMap[AutumnYellow] + "   Default Fetch")
	fmt.Println(ColorMap[AutumnGreen] + "Flags: " + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -h or --help:" + ColorMap[AutumnYellow] + "      Prints help menu" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -v or --version:" + ColorMap[AutumnYellow] + "   Prints version of Edison" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -pa or --print-ascii:" + ColorMap[AutumnYellow] + "     Shows ASCII in config.json" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -gv or --gitea-version:" + ColorMap[AutumnYellow] + "     Prints current version of Gitea" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -pc or --print-config:" + ColorMap[AutumnYellow] + "    Prints current configuration file." + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -gu or --get-user:" + ColorMap[AutumnYellow] + "      Prints current Gitea user" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -lr or --list-repos:" + ColorMap[AutumnYellow] + "     Prints repos and ssh clone links" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -lor or --list-org-repos:" + ColorMap[AutumnYellow] + "      Prints organization repos and ssh links" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -rc or --repo-count:" + ColorMap[AutumnYellow] + "      Prints the number of repos" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -cr or --create-repo:" + ColorMap[AutumnYellow] + "      Starts prompt to create new repository" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -pul or --print-user-langs:" + ColorMap[AutumnYellow] + "     Prints percent of languages used" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -tul or --top-user-langs:" + ColorMap[AutumnYellow] + "       Prints a top used languages" + ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen] + "   -loc or --lines-of-code:" + ColorMap[AutumnYellow] + "   Prints LOC across all repos" + ColorMap[DefaultReset])
}

func print_ascii_art(ascii_config []string) {
	for _, line := range ascii_config {
		fmt.Println(ColorMap[AutumnGreen] + line + ColorMap[DefaultReset])
	}
}

func print_edison_version() {
	fmt.Println(ColorMap[AutumnGreen] + "Edison Version:" + ColorMap[AutumnYellow] + " 1.0.0" + ColorMap[DefaultReset])

}

func make_config_dir() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configDir := filepath.Join(homeDir, ".config", "edison")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	configFile := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig := map[string]any{
			"token": "gitea_token",
			"URI":   "gitea_url",
			"ascii_art": []string{
				" ____  ____  ____  ____  ____  ____ ",
				"||E ||||d ||||i ||||s ||||o ||||n ||",
				"||__||||__||||__||||__||||__||||__||",
				"|/__\\||/__\\||/__\\||/__\\||/__\\||/__\\|",
			},
			"top_langs":      3,
			"repo_count":     true,
			"include_orgs":   false,
			"gitea_user":     true,
			"gitea_version":  true,
			"edison_version": false,
		}

		file, err := os.Create(configFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(defaultConfig)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Created default config.json")
	} else {
		//config exists
	}
}

func load_config() Configuration {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configFile := filepath.Join(homeDir, ".config", "edison", "config.json")

	jsonFile, err := os.Open(configFile)
	if err != nil {
		log.Println("Error opening config file:", err)
		return Configuration{}
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println("Error reading config file:", err)
		return Configuration{}
	}

	var user_configuration Configuration
	err = json.Unmarshal(byteValue, &user_configuration)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return Configuration{}
	}
	return user_configuration
}

func print_config(user_configuration Configuration) {
	fmt.Println(ColorMap[AutumnGreen]+"Token: "+ColorMap[AutumnYellow], user_configuration.Token, ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen]+"URI: "+ColorMap[AutumnYellow], user_configuration.URI, ColorMap[DefaultReset])

	fmt.Println(ColorMap[AutumnGreen] + "ASCII Art:")
	for _, line := range user_configuration.Ascii_Art {
		fmt.Println(ColorMap[AutumnGreen] + line + ColorMap[DefaultReset])
	}
	fmt.Print(ColorMap[DefaultReset])

	fmt.Println(ColorMap[AutumnGreen]+"Top Langs: "+ColorMap[AutumnYellow], user_configuration.Top_Langs, ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen]+"Repo Count: "+ColorMap[AutumnYellow], user_configuration.Repo_Count, ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen]+"Include Orgs: "+ColorMap[AutumnYellow], user_configuration.Include_Orgs, ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen]+"Gitea User: "+ColorMap[AutumnYellow], user_configuration.Gitea_User, ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen]+"Gitea Version: "+ColorMap[AutumnYellow], user_configuration.Gitea_Version, ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen]+"Edison Version: "+ColorMap[AutumnYellow], user_configuration.Edison_Version, ColorMap[DefaultReset])
}

func create_gitea_client(config_uri string, config_token string) *gitea.Client {
	client, err := gitea.NewClient(config_uri, gitea.SetToken(config_token))
	if err != nil {
		log.Println(ColorMap[AutumnRed]+"Error creating Gitea client:", err, ColorMap[DefaultReset])
	}
	return client
}

func print_gitea_server_version(client *gitea.Client) {
	version, _, err := client.ServerVersion()
	if err != nil {
		log.Println("Error getting Gitea server version:", err)
	}
	fmt.Println(ColorMap[AutumnGreen]+"Gitea Version:"+ColorMap[AutumnYellow], version, ColorMap[DefaultReset])
}

func print_gitea_user(client *gitea.Client) {
	gitea_user_data, _, err := client.GetMyUserInfo()
	if err != nil {
		log.Println("Error getting Gitea user:", err)
	}
	fmt.Println(ColorMap[AutumnGreen]+"Username:"+ColorMap[AutumnYellow], gitea_user_data.UserName, ColorMap[DefaultReset])
}

func print_user_repos(client *gitea.Client) {
	user_repos, _, err := client.ListMyRepos(gitea.ListReposOptions{})
	if err != nil {
		log.Println("Error getting Gitea user repos:", err)
	}
	for _, repo := range user_repos {
		gitea_user_data, _, err := client.GetMyUserInfo()
		if err != nil {
			log.Println("Error getting Gitea user:", err)
		}
		if repo.Owner.UserName == gitea_user_data.UserName {
			fmt.Println(ColorMap[AutumnGreen], repo.Name, ColorMap[DefaultReset])
			fmt.Println(ColorMap[AutumnYellow], repo.SSHURL, ColorMap[DefaultReset])
		}
	}
}

func print_org_repos(client *gitea.Client) {
	user_orgs, _, err := client.ListMyOrgs(gitea.ListOrgsOptions{})
	if err != nil {
		log.Println("Error getting Gitea user orgs:", err)
	}
	for _, org := range user_orgs {
		var org_name = org.UserName
		fmt.Println(ColorMap[AutumnGreen], org.UserName, ColorMap[DefaultReset])
		org_repos, _, err := client.ListOrgRepos(org_name, gitea.ListOrgReposOptions{})
		if err != nil {
			log.Println("Error getting Gitea org repos:", err)
		}
		for _, repo := range org_repos {
			fmt.Println(ColorMap[AutumnGreen], repo.Name, ColorMap[DefaultReset])
			fmt.Println(ColorMap[AutumnYellow], repo.SSHURL, ColorMap[DefaultReset])
		}
	}
}

func create_gitea_repo(client *gitea.Client) {
	var user_options = gitea.CreateRepoOption{}

	fmt.Print(ColorMap[AutumnGreen]+"Repo Name: ", ColorMap[AutumnYellow])
	user_options.Name, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	user_options.Name = strings.TrimSpace(user_options.Name)

	fmt.Print(ColorMap[AutumnGreen]+"Repo Description: ", ColorMap[AutumnYellow])
	user_options.Description, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	user_options.Description = strings.TrimSpace(user_options.Description)

	var user_private_response string
	fmt.Print(ColorMap[AutumnGreen]+"Private Repo (y/n): ", ColorMap[AutumnYellow])
	fmt.Scanln(&user_private_response)
	fmt.Println(ColorMap[DefaultReset])

	user_private_response = strings.ToLower(user_private_response)

	if user_private_response == "y" {
		fmt.Println(ColorMap[AutumnGreen]+"The repo will be private.", ColorMap[DefaultReset])
		user_options.Private = true
	} else if user_private_response == "n" {
		fmt.Println(ColorMap[AutumnGreen]+"The repo will NOT be private.", ColorMap[DefaultReset])
		user_options.Private = false
	} else {
		fmt.Println(ColorMap[AutumnGreen]+"The repo will be private.", ColorMap[DefaultReset])
		user_options.Private = true
	}

	repo_creation, _, err := client.CreateRepo(user_options)
	if err != nil {
		log.Println("Error create Gitea repos:", err)
	}
	fmt.Println(ColorMap[AutumnGreen]+"Created:"+ColorMap[AutumnYellow], repo_creation.Name, ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen]+"Description:"+ColorMap[AutumnYellow], repo_creation.Description, ColorMap[DefaultReset])
	fmt.Println(ColorMap[AutumnGreen]+"Private:"+ColorMap[AutumnYellow], repo_creation.Private, ColorMap[DefaultReset])
}

func print_user_repo_count(config_repo_count bool, config_include_orgs bool, client *gitea.Client) {
	var user_repo_count = 0
	if config_repo_count == false {
		return
	}

	user_repos, _, err := client.ListMyRepos(gitea.ListReposOptions{})
	if err != nil {
		log.Println("Error getting Gitea user repos:", err)
	}

	if config_include_orgs == false {
		gitea_user_data, _, err := client.GetMyUserInfo()
		if err != nil {
			log.Println("Error getting Gitea user:", err)
		}

		for _, repo := range user_repos {
			if repo.Owner.UserName == gitea_user_data.UserName {
				user_repo_count++
			}
		}
	} else {
		user_repo_count = len(user_repos)
	}
	fmt.Println(ColorMap[AutumnGreen]+"Repositories:"+ColorMap[AutumnYellow], user_repo_count, ColorMap[DefaultReset])
}

func print_user_langs(client *gitea.Client, include_org bool, top_only bool, user_top_langs int) {
	langStats := make(map[string]int64)
	var totalLines int64

	user_repos, _, err := client.ListMyRepos(gitea.ListReposOptions{})
	if err != nil {
		log.Println("Error getting Gitea user repos:", err)
		return
	}

	gitea_user_data, _, err := client.GetMyUserInfo()
	if err != nil {
		log.Println("Error getting Gitea user:", err)
	}

	if !include_org {
		for _, repo := range user_repos {
			if repo.Owner.UserName == gitea_user_data.UserName {
				repo_lang_data, _, err := client.GetRepoLanguages(gitea_user_data.UserName, repo.Name)
				if err != nil {
					log.Println("Error getting repo lang data:", err)
					continue
				}

				for lang, size := range repo_lang_data {
					langStats[lang] += size
					totalLines += size
				}
			}
		}

	} else {
		for _, repo := range user_repos {
			repo_lang_data, _, err := client.GetRepoLanguages(gitea_user_data.UserName, repo.Name)
			if err != nil {
				log.Println("Error getting repo lang data:", err)
				continue
			}
			for lang, size := range repo_lang_data {
				langStats[lang] += size
				totalLines += size
			}
		}
	}
	var sortedLangs []langStat
	for lang, size := range langStats {
		percentage := (float64(size) / float64(totalLines)) * 100
		sortedLangs = append(sortedLangs, langStat{Lang: lang, Percentage: percentage})
	}

	sort.Slice(sortedLangs, func(i, j int) bool {
		return sortedLangs[i].Percentage > sortedLangs[j].Percentage
	})

	if !top_only {
		for _, entry := range sortedLangs {
			fmt.Printf("%s%s: %s%.2f%%%s\n", ColorMap[AutumnGreen], entry.Lang, ColorMap[AutumnYellow], entry.Percentage, ColorMap[DefaultReset])
		}
	} else {
		fmt.Println(ColorMap[AutumnGreen] + "Top Languages:" + ColorMap[DefaultReset])
		for i, entry := range sortedLangs {
			if i >= user_top_langs {
				break
			}
			fmt.Printf("%s%s:%s %.2f%%%s\n", ColorMap[AutumnGreen], entry.Lang, ColorMap[AutumnYellow], entry.Percentage, ColorMap[DefaultReset])
		}
	}
}

func print_user_total_loc(client *gitea.Client, include_org bool) {
	var totalLines int64

	user_repos, _, err := client.ListMyRepos(gitea.ListReposOptions{})
	if err != nil {
		log.Println("Error getting Gitea user repos:", err)
		return
	}

	gitea_user_data, _, err := client.GetMyUserInfo()
	if err != nil {
		log.Println("Error getting Gitea user:", err)
	}

	if !include_org {
		for _, repo := range user_repos {
			if repo.Owner.UserName == gitea_user_data.UserName {
				repo_lang_data, _, err := client.GetRepoLanguages(gitea_user_data.UserName, repo.Name)
				if err != nil {
					log.Println("Error getting repo lang data:", err)
					continue
				}
				for _, size := range repo_lang_data {
					totalLines += size
				}
			}
		}
	} else {
		for _, repo := range user_repos {
			repo_lang_data, _, err := client.GetRepoLanguages(gitea_user_data.UserName, repo.Name)
			if err != nil {
				log.Println("Error getting repo lang data:", err)
				continue
			}
			for _, size := range repo_lang_data {
				totalLines += size
			}
		}
	}
	fmt.Println(ColorMap[AutumnGreen]+"LOC:"+ColorMap[AutumnYellow], totalLines, ColorMap[DefaultReset])
}
