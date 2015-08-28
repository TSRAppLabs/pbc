package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"stash.tsrapplabs.com/ut/pbc"
)

func getSignConfig() (pbc.SignConfig, error) {
	return pbc.SignConfig{
		Cert:   cert,
		Signer: signer,
		Key:    key,
		Pass:   pass,
	}, nil
}

func main() {

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

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

var cert string
var signer string
var key string
var pass string

var nosign bool
var nozip bool

var help bool

func init() {
	flag.StringVar(&cert, "cert", "", "Cert with which to sign passbook")
	flag.StringVar(&signer, "signer", "", "Signing identity")
	flag.StringVar(&key, "key", "", "Key")
	flag.StringVar(&pass, "pass", "", "Pass")

	flag.BoolVar(&nosign, "no-sign", false, "Will not sign pass")
	flag.BoolVar(&nozip, "no-zip", false, "Will not zip")

	flag.BoolVar(&help, "help", false, "Help command")
}
