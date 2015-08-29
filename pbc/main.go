package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"stash.tsrapplabs.com/ut/pbc"
)

func main() {
	pbc.InitDataDir()
	rootCmd.Execute()
}

var passname string

var nosign bool
var nozip bool

var help bool

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use: "pbc",
	}

	rootCmd.AddCommand(mkBuildCommand())
	rootCmd.AddCommand(mkProfileCommand())
}

func mkBuildCommand() *cobra.Command {
	var profilename string

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "builds a pkpass",
		Long:  "builds a pkpass",
		Run: func(cmd *cobra.Command, args []string) {
			root := "."
			if len(args) > 0 {
				root = args[0]
			}

			profile, err := pbc.GetProfile(profilename)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v", err.Error())
				return
			}

			file, err := os.Create(passname)

			if err != nil {
				fmt.Println(err)
				log.Fatal(err)
			}

			err = pbc.Compile(root, profile, file)

			if err != nil {
				fmt.Println(err.Error())
				log.Fatal(err)

			}
		},
	}

	buildCmd.Flags().StringVarP(&profilename, "profile", "p", "", "Profile to use")
	buildCmd.Flags().StringVarP(&passname, "name", "n", "pass.pkpass", "Resulting passbook file")

	return buildCmd
}

func mkProfileCommand() *cobra.Command {
	var name string
	var p12path string

	addProfile := func(cmd *cobra.Command, args []string) {
		if name != "" && p12path != "" {
			profile, err := pbc.CreateProfile(name, p12path)

			if err != nil {
				log.Fatal(err)
			}

			pbc.AddProfile(profile)
		}
	}

	profCmd := &cobra.Command{
		Use:   "profile",
		Short: "profile management",
		Long:  "manages a profile to add, rm, ls",
	}

	profAddCmd := &cobra.Command{
		Use:   "add",
		Short: "adds a profile",
		Run:   addProfile,
	}
	profAddCmd.Flags().StringVarP(&name, "profile-name", "p", "", "Name to give the profile")
	profAddCmd.Flags().StringVarP(&p12path, "cert", "c", "", "Cert to create profile with")

	profLsCmd := &cobra.Command{
		Use:   "ls",
		Short: "lists profiles",
		Long:  "lists profiles",
		Run: func(cmd *cobra.Command, args []string) {
			for _, prof := range pbc.ListProfiles() {
				fmt.Printf("\t%v\n", prof.Name)
			}
		},
	}

	profRmCmd := &cobra.Command{
		Use:   "rm",
		Short: "removes profiles",
		Long:  "removes all profile specified",
		Run: func(cmd *cobra.Command, args []string) {

			for _, arg := range args {
				pbc.DelProfile(arg)
			}

		},
	}

	profCmd.AddCommand(profAddCmd)
	profCmd.AddCommand(profLsCmd)
	profCmd.AddCommand(profRmCmd)

	return profCmd
}
