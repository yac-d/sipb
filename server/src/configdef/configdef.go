package configdef

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
)

type Config struct {
	BinPath    string `yaml:"BinPath"`
	BinDir     string
	WebpageDir string `yaml:"WebpageDir"`
}

//ReadFromYAML reads config information from the YAML file at the specified path
func (c *Config) ReadFromYAML(fp string) {
	configfileBytes, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatalln("Unable to read config file")
	}
	yaml.Unmarshal(configfileBytes, c)
	c.BinDir = path.Join(c.WebpageDir, c.BinPath)
}

