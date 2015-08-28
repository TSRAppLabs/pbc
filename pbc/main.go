package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"stash.tsrapplabs.com/ut/pbc"
)

func main() {
	rootCmd.Execute()
}

var cert string
var signer string
var key string
var pass string

var nosign bool
var nozip bool

var help bool

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use: "pbc",
	}

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "builds a pkpass",
		Long:  "builds a pkpass",
		Run:   buildCommandRun,
	}

	buildCmd.Flags().StringVarP(&cert, "cert", "c", "", "Certificate to sign")
	buildCmd.Flags().StringVarP(&signer, "signer", "s", "", "Certificate signing the certificate")
	buildCmd.Flags().StringVarP(&key, "key", "k", "", "Key")
	buildCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for certificate")

	rootCmd.AddCommand(buildCmd)
}

func buildCommandRun(cmd *cobra.Command, args []string) {

	fmt.Println(args)

	root := "."

	file, err := os.Create("vr.pkpass")

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

	compileOptions := pbc.CompileOptions{
		NoZip:       nozip,
		NoSignature: nosign,
	}

	err = pbc.Compile(root, config, compileOptions, file)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
}
