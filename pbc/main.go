package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"stash.tsrapplabs.com/ut/pbc"
)

func getSignConfig() (pbc.SignConfig, error) {
	var signConfig pbc.SignConfig
	if cert == "" {
		return signConfig, errors.New("Please specify a cert to sign the passbook with")
	}

	if signer == "" {
		return signConfig, errors.New("Please specify a signer for the passbook")
	}

	if key == "" {
		return signConfig, errors.New("Please specify a intermediate Key")
	}

	if pass == "" {
		return signConfig, errors.New("Please supply a password")
	}
	return pbc.SignConfig{
		Cert:   cert,
		Signer: signer,
		Key:    key,
		Pass:   pass,
	}, nil
}

func main() {

	flag.Parse()

	config, err := getSignConfig()

	root := "."
	if len(flag.Args()) > 1 {
		root = flag.Args()[1]
	}

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	file, err := os.Create("vr.pkpass")

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	err = pbc.Compile(root, config, file)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
}

var cert string
var signer string
var key string
var pass string

func init() {
	flag.StringVar(&cert, "cert", "", "Cert with which to sign passbook")
	flag.StringVar(&signer, "signer", "", "Signing identity")
	flag.StringVar(&key, "key", "", "Key")
	flag.StringVar(&pass, "pass", "", "Pass")
}
