package configdef

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"os"
	"strconv"
)

type Config struct {
	BinPath    string `yaml:"BinPath"`
	BinDir     string
	WebpageDir string `yaml:"WebpageDir"`
	Port       int    `yaml:"Port"`
	BindAddr   string `yaml:"BindAddr"`
	MaxFileCnt int    `yaml:"MaxFileCnt"`
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

//ReadFromEnvVars reads config information from environment variables
func (c *Config) ReadFromEnvVars() {
	if val, set := os.LookupEnv("SIPB_PORT"); set {
		p, _ := strconv.Atoi(val)
		c.Port = p
	}
	if val, set := os.LookupEnv("SIPB_BIN_PATH"); set {
		c.BinPath = val
		c.BinDir = path.Join(c.WebpageDir, c.BinPath)
	}
	if val, set := os.LookupEnv("SIPB_WEBPAGE_DIR"); set {
		c.WebpageDir = val
	}
	if val, set := os.LookupEnv("SIPB_BIND_ADDR"); set {
		c.BindAddr = val
	}
	if val, set := os.LookupEnv("SIPB_MAX_FILE_CNT"); set {
		cnt, _ := strconv.Atoi(val)
		c.MaxFileCnt = cnt
	}
}
