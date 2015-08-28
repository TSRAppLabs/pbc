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
	var cert string
	var signer string
	var key string
	var pass string

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "builds a pkpass",
		Long:  "builds a pkpass",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)

			root := "."

			file, err := os.Create(passname)

			if err != nil {
				fmt.Println(err)
				log.Fatal(err)
			}

			config := pbc.SignConfig{
				Cert:   cert,
				Signer: signer,
				Key:    key,
				Pass:   pass,
			}

			err = pbc.Compile(root, config, file)

			if err != nil {
				fmt.Println(err.Error())
				log.Fatal(err)
			}
		},
	}

	buildCmd.Flags().StringVarP(&cert, "cert", "c", "", "Certificate to sign")
	buildCmd.Flags().StringVarP(&signer, "signer", "s", "", "Certificate signing the certificate")
	buildCmd.Flags().StringVarP(&key, "key", "k", "", "Key")
	buildCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for certificate")
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
