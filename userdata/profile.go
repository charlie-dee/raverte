package userdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"raverte/appdata"
)

type Profile struct {
	Name     string `json:"name"`
	Keystore bool   `json:"keystore"`
}

func (p *Profile) ReturnSelf() *Profile {
	return p
}

// Sets default values for a fresh Raverte installation.
//
// Returns error if: profile already exists, can not create profile or can not write to profile.json.
func (p *Profile) InitialiseProfile() error {
	filePath, err := GetRaverteAsset(appdata.PROFILE)
	if err != nil {
		return err
	}
	if err := checkRaverteAsset(filePath); err != nil {
		if err.Error() == "file does not exist" {
			if err = configureRaverteAsset(filePath, appdata.PROFILE); err != nil {
				return err
			}
		}
	} else {
		return errors.New("profile already exists")
	}

	p.setDefaultValues()
	if err := p.writeProfile(); err != nil {
		return err
	}
	return nil
}

func (p *Profile) LoadProfile() error {
	filePath, err := GetRaverteAsset(appdata.PROFILE)
	if err != nil {
		return err
	}

	if err = checkRaverteAsset(filePath); err != nil {
		return err
	}

	readData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading from profile: %s", err.Error())
	}

	if err = json.Unmarshal(readData, p); err != nil {
		return fmt.Errorf("unable to parse data from profile: %s", err.Error())
	}
	return nil
}

// Sets keystore value and commits to profile.json
//
// Returns error if: invalid profile(does not exist/permissions have been changed) or is unable to write to profile.json.
func (p *Profile) UpdateKeystore(value bool) error {
	p.Keystore = value

	if err := p.writeProfile(); err != nil {
		return err
	}

	return nil
}

// Sets default values of profile.json
func (p *Profile) setDefaultValues() {
	p.Name = "Trader J"
	p.Keystore = false
}

// Commits Profile attributes that are currently stored in memory to profile.json
//
// Returns error if: invalid keystore (does not exist/permissions have been changed) or unable to write to profile.json
func (p *Profile) writeProfile() error {
	profilePath, err := GetRaverteAsset(appdata.PROFILE)
	if err != nil {
		return err
	}

	newProfile, err := json.Marshal(p)
	if err != nil {
		return err
	}

	if err := os.WriteFile(profilePath, newProfile, 0); err != nil {
		return err
	}
	return nil
}
