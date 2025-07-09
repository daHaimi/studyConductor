package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var Cfg *Config

type BlockDevice struct {
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	Mountpoint string        `json:"mountpoint	"`
	Uuid       string        `json:"uuid"`
	Children   []BlockDevice `json:"children"`
}

func GetStorage(config *Config) (string, error) {
	Cfg = config
	path := config.Study.Storage.GetString("path")
	if Cfg.Study.Storage["local"] != "true" {
		// Use USB Drive
		drive, err := Detect()
		if err != nil {
			return "", err
		}
		if drive == nil {
			return "", fmt.Errorf("no USB drive found within specifications")
		}

		if drive.Mountpoint == "" {
			return "", fmt.Errorf("device %s is not mounted. Please mount and try again", drive.Name)
		}
		path = filepath.Join(drive.Mountpoint, path)
	}
	// Check if path already exists; if not: create
	fmt.Println("Your study data will be stored at: ", path)
	return path, os.MkdirAll(path, 0777)
}

// Detect returns a list of file paths pointing to the root folder of
// USB storage devices connected to the system.
func Detect() (*BlockDevice, error) {
	var drives map[string][]BlockDevice

	out, err := exec.Command("lsblk", "-J", "-o", "NAME,TYPE,MOUNTPOINT,UUID").CombinedOutput()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &drives)
	if err != nil {
		return nil, err
	}

	for _, bd := range drives["blockdevices"] {
		for _, part := range bd.Children {
			if part.Uuid == Cfg.Study.Storage["id"].(string) {
				return &part, nil
			}
		}
	}
	return nil, nil
}
