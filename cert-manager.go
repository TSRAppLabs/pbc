package pbc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Profile struct {
	Name string `json:"name"`
}

func (profile Profile) getBasePath() string {
	return filepath.Join(getDataDir(), "profiles", profile.Name)
}

func (profile Profile) getKeyPath() string {
	return filepath.Join(profile.getBasePath(), "key.pem")
}

func (profile Profile) getCertPath() string {
	return filepath.Join(profile.getBasePath(), "cert.pem")
}

func ListProfiles() []Profile {
	indexPath := filepath.Join(getDataDir(), "index.json")
	file, err := os.Open(indexPath)

	if err != nil {
		file, err := os.Create(indexPath)
		if err != nil {
			return []Profile{}
		}

		json.NewEncoder(file).Encode([]Profile{})
		return []Profile{}
	}

	data := []Profile{}
	err = json.NewDecoder(file).Decode(&data)

	if err != nil {
		return []Profile{}
	}

	return data
}

func SaveProfiles(profiles []Profile) {
	file, err := os.Create(filepath.Join(getDataDir(), "index.json"))

	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(file).Encode(profiles)
}

func AddProfile(profile Profile) {
	SaveProfiles(append(ListProfiles(), profile))
}

func DelProfile(name string) {
	profiles := ListProfiles()

	mfd := -1
	for i, p := range profiles {
		if p.Name == name {
			mfd = i
		}
	}

	if mfd != -1 {
		os.RemoveAll(profiles[mfd].getBasePath())
		SaveProfiles(append(profiles[:mfd], profiles[mfd+1:]...))
	}
}

func GetProfile(name string) (Profile, error) {
	profiles := ListProfiles()

	for _, prof := range profiles {
		if prof.Name == name {
			return prof, nil
		}
	}

	return Profile{}, errors.New("Profile not found")
}

func CreateProfile(name, p12path string) (Profile, error) {
	for _, p := range ListProfiles() {
		if p.Name == name {
			return Profile{}, errors.New("Profile exists")
		}
	}

	profile := Profile{
		Name: name,
	}

	dir := filepath.Join(getDataDir(), "profiles", profile.Name)
	os.MkdirAll(dir, 0755)
	os.Chmod(dir, 0755)

	createCertKey(p12path, profile.getCertPath(), profile.getKeyPath())
	return profile, nil
}

func copyWWDR() error {
	if _, err := os.Stat(filepath.Join(getDataDir(), "wwdr.pem")); err != nil {
		wwdr, err := os.Open(ExpandPath("${GOPATH}/src/stash.tsrapplabs.com/ut/pbc/wwdr.pem"))

		if err != nil {
			return err
		}

		defer wwdr.Close()

		file, err := os.Create(filepath.Join(getDataDir(), "wwdr.pem"))

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(file, wwdr)

		return err
	}

	return nil
}

func InitDataDir() {
	if !exists(getDataDir()) {
		fmt.Println("Adding datadir")
		os.MkdirAll(getDataDir(), 0750)
		os.Chmod(getDataDir(), 0750)
	}

	if !exists(filepath.Join(getDataDir(), "index.json")) {
		fmt.Println("Adding index.json")
		file, err := os.Create(filepath.Join(getDataDir(), "index.json"))
		if err == nil {
			_, err := file.WriteString("{}")

			if err != nil {
				fmt.Println(err)
			}
			file.Close()
		}
	}

	if !exists(filepath.Join(getDataDir(), "wwdr.pem")) {
		fmt.Println("Copying wwdr.pem")
		if err := copyWWDR(); err != nil {
			fmt.Println(err)
		}
	}

	if !exists(filepath.Join(getDataDir(), "profiles")) {
		fmt.Println("Adding profiles")
		os.MkdirAll(filepath.Join(getDataDir(), "profiles"), 0750)
		os.Chmod(getDataDir(), 0750)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
