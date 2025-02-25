package main

import "code.gitea.io/sdk/gitea"
import "fmt"

func main() {
	client, err := gitea.NewClient("http://192.168.7.2:3029", gitea.SetToken("4134b116be73b40c8bc8051dd29fc76f64d53f23"))
	if err != nil {
		fmt.Println("Error creating Gitea client:", err)
		return
	}

	user, _, err := client.GetUserInfo("c1uckie")
	if err != nil {
		fmt.Println("Error fetching user:", err)
		return
	}

	// Print the user info
	fmt.Println("User:", user.UserName)

}
