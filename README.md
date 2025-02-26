# edison

Gitea CLI fetch tool written in Go

## Purpose
Just a small CLI I use to keep my workflow inside the CLI.

## Project Setup
```sh
go mod init gitea.ohara.local/c1uckie/edison
go get code.gitea.io/sdk/gitea
```

## Installation
To be determined

## Usage
1. Create an API token
2. Add the Gitea API token and URI to your config.json
3. run `edison`

## Configuration
Default JSON configuration
```json
{
  "token": "8465cf9a4321b8c6b4a45d376bf13bdca2a854df",
  "URI": "http://192.168.7.2:3029",
  "ascii_art": [
    " ____  ____  ____  ____  ____  ____ ",
    "||E ||||d ||||i ||||s ||||o ||||n ||",
    "||__||||__||||__||||__||||__||||__||",
    "|/__\\||/__\\||/__\\||/__\\||/__\\||/__\\|"
  ],
  "top_langs": 3,
  "repo_count": true,
  "include_orgs": false,
  "gitea_user": true,
  "gitea_version": true,
  "edison_version": true
}
```
