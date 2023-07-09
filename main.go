package main

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

const Version = "v0.0.1"

var rootCmd = &cobra.Command{
	Use:   "git-profile",
	Short: "Helper for managing multiple git profiles for the same user",
	Long:  `Complete documentation is available at https://github.com/gigurra/flycd`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	Run: func(cmd *cobra.Command, args []string) {

		users := getAvailableProfiles()

		fmt.Printf("Available profiles: \n")
		for profileName, v := range users {
			fmt.Printf(" - %s: %s, %s, %s\n", profileName, v.SSHConfig.IdentityFile, v.GitConfig.UserName, v.GitConfig.UserEmail)
		}
	},
}

func getAvailableProfiles() Users {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home dir: %v\n", err)
		os.Exit(1)
	}

	profilesFilePath := homeDir + "/.ssh/git-profiles.json"

	// list profiles from ~/.ssh/git-profiles.json
	contents, err := os.ReadFile(profilesFilePath)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", profilesFilePath, err)
		os.Exit(1)
	}

	// parse json
	var users Users
	err = json.Unmarshal(contents, &users)
	if err != nil {
		fmt.Printf("Error parsing %s: %v\n", profilesFilePath, err)
		os.Exit(1)
	}

	return users
}

var setCmd = &cobra.Command{
	Use:   "set <profile name>",
	Short: "Set the git profile (ssh config and git config)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		profileName := args[0]

		users := getAvailableProfiles()

		fmt.Printf("Available profiles: \n")
		for profileName, v := range users {
			fmt.Printf(" - %s: %s, %s, %s\n", profileName, v.SSHConfig.IdentityFile, v.GitConfig.UserName, v.GitConfig.UserEmail)
		}

		fmt.Printf("Setting profile '%s'\n", profileName)
		profile, found := lo.FindKeyBy(users, func(key string, value User) bool {
			return key == profileName
		})

		if !found {
			fmt.Printf("Error: profile '%s' not found\n", profileName)
			os.Exit(1)
		}

		// set ssh config
		print(profile)

		fmt.Printf("Not implemented yet, sorry :(\n")
		os.Exit(1)
	},
}

func main() {

	// Check that required applications are installed
	requiredApps := []string{"git", "ssh"}
	for _, app := range requiredApps {
		_, err := exec.LookPath(app)
		if err != nil {
			fmt.Printf("Error: required app '%s' not found in PATH\n", app)
			os.Exit(1)
		}
	}

	// prepare cli
	rootCmd.AddCommand(listCmd, setCmd)

	// run cli
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type SSHConfig struct {
	Host           string `json:"Host"`
	User           string `json:"User"`
	IdentityFile   string `json:"IdentityFile"`
	IdentitiesOnly string `json:"IdentitiesOnly"`
	IdentityAgent  string `json:"IdentityAgent"`
}

type GitConfig struct {
	UserName  string `json:"user.name"`
	UserEmail string `json:"user.email"`
}

type User struct {
	SSHConfig SSHConfig `json:"ssh_config"`
	GitConfig GitConfig `json:"git_config"`
}

type Users map[string]User
