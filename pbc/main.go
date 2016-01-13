package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tsrapplabs/pbc"
)

func main() {
	pbc.InitDataDir()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading config: %v", err)
	}

	rootCmd.Execute()
}

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use: "pbc",
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(os.ExpandEnv("${HOME}/.pbc"))
	viper.AddConfigPath("/etc/pbc")

	viper.SetDefault("core.datadir", "${HOME}/.pbc")

	rootCmd.AddCommand(mkBuildCommand())
	rootCmd.AddCommand(mkProfileCommand())
	rootCmd.AddCommand(mkLintCommand())
}

func mkBuildCommand() *cobra.Command {
	var profile string
	var name string

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "builds a pkpass",
		Long:  "builds a pkpass",
		Run: func(cmd *cobra.Command, args []string) {
			root := "."
			if len(args) > 0 {
				root = args[0]
			}
			if profile != "" {
				viper.SetDefault("build.profile", profile)
			}

			if name != "" {
				viper.SetDefault("build.name", name)
			}

			profile, err := pbc.GetProfile(viper.GetString("build.profile"))

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
				return
			}

			file, err := os.Create(viper.GetString("build.name"))

			if err != nil {

				fmt.Printf("Trying to create file: %v, %v", viper.GetString("build.name"), err)
				os.Exit(1)
			}

			err = pbc.Compile(root, profile, file)

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)

			}
		},
	}
	buildCmd.Flags().StringVarP(&profile, "profile", "p", "", "Profile to use")
	buildCmd.Flags().StringVarP(&name, "name", "n", "", "Resulting passbook file")
	viper.SetDefault("build.name", "pass.pkpass")
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
	profAddCmd.Flags().StringVarP(&name, "profile", "p", "", "Name to give the profile")
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

func mkLintCommand() *cobra.Command {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "checks a pass for mistakes",
		Long:  "checks a pass for mistakes",
		Run: func(cmd *cobra.Command, args []string) {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			warn, err := pbc.LintPass(dir)

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			for _, msg := range warn {
				fmt.Println(msg)
			}

		},
	}

	return lintCmd
}
