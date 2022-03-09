package configdef

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	BinDir     string `yaml:"BinDir"`
	WebpageDir string `yaml:"WebpageDir"`
}

//ReadFromYAML reads config information from the YAML file at the specified path
func (c *Config) ReadFromYAML(fp string) {
	configfileBytes, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatalln("Unable to read config file")
	}
	yaml.Unmarshal(configfileBytes, c)
}

