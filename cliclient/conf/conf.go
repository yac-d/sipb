package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const SBCRC = "sbcrc"

type Config struct {
	Server string `json:"Server"`
}

func filepath() string {
	fp, _ := os.UserConfigDir()
	return path.Join(fp, SBCRC)
}

// New creates a new config, prompting for new
// values and persisting if required
func New() (*Config, error) {
	c := &Config{}

	if _, err := os.Stat(filepath()); os.IsNotExist(err) {
		return c, c.PromptAndPersist()
	}

	bytes, err := ioutil.ReadFile(filepath())
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, c)
	return c, err
}

func (c *Config) PromptAndPersist() error {
	fmt.Println("--- Configuration ---")
	fmt.Println("Enter pastebin server (Ex: https://user:pass@example.com:1234/):")
	fmt.Scanf("%s", &c.Server)

	bytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	file, err := os.Create(filepath())
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(bytes)
	return err
}
