package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

const Version = "v0.0.3"

var rootCmd = &cobra.Command{
	Use:   "git-profile [profile-name]",
	Short: "Helper for managing multiple git profiles for the same user",
	Run: func(cmd *cobra.Command, args []string) {

		users := getAvailableProfiles()

		if len(args) == 0 {

			fmt.Printf("Available profiles: \n")
			for profileName, v := range users {
				fmt.Printf(" - %s: %s, %s, %s\n", profileName, v.SSHConfig.IdentityFile, v.GitConfig.UserName, v.GitConfig.UserEmail)
			}
		} else {

			profileName := args[0]

			users := getAvailableProfiles()

			fmt.Printf("Setting profile '%s'\n", profileName)
			profile, found := users[profileName]

			if !found {
				fmt.Printf("Error: profile '%s' not found\n", profileName)
				os.Exit(1)
			}

			// set ssh config
			fmt.Printf("Setting Git parameters\n")
			fmt.Printf(" - git config --global user.name %s\n", profile.GitConfig.UserName)
			err := exec.Command("git", "config", "--global", "user.name", profile.GitConfig.UserName).Run()
			if err != nil {
				fmt.Printf("Error setting git config user.name: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf(" - git config --global user.email %s\n", profile.GitConfig.UserEmail)
			err = exec.Command("git", "config", "--global", "user.email", profile.GitConfig.UserEmail).Run()
			if err != nil {
				fmt.Printf("Error setting git config user.email: %v\n", err)
				os.Exit(1)
			}

			// set ssh config
			fmt.Printf("Setting SSH parameters\n")

			newContents := fmt.Sprintf(
				`Host *
    User "git"
    IdentityFile "%s"
    IdentitiesOnly "yes"
    IdentityAgent "%s"
`,
				profile.SSHConfig.IdentityFile,
				profile.SSHConfig.IdentityAgent,
			)

			fmt.Printf("New ssh config:\n")
			print(newContents)

			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("Error getting home dir: %v\n", err)
				os.Exit(1)
			}

			sshConfigFilePath := homeDir + "/.ssh/config"
			err = os.WriteFile(sshConfigFilePath, []byte(newContents), 0644)
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", sshConfigFilePath, err)
				os.Exit(1)
			}
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

	// run cli
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type SSHConfig struct {
	IdentityFile  string `json:"IdentityFile"`
	IdentityAgent string `json:"IdentityAgent"`
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
