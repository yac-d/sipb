package configdef

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	BinPath     string `yaml:"BinPath" env:"SIPB_BIN_PATH"`
	BinDir      string
	WebpageDir  string `yaml:"WebpageDir" env:"SIPB_WEBPAGE_DIR"`
	Port        int    `yaml:"Port" env:"SIPB_PORT"`
	BindAddr    string `yaml:"BindAddr" env:"SIPB_BIND_ADDR"`
	MaxFileCnt  int    `yaml:"MaxFileCnt" env:"SIPB_MAX_FILE_CNT"`
	MaxFileSize int64  `yaml:"MaxFileSize" env:"SIPB_MAX_FILE_SIZE"` // Bytes
}

//ReadFromYAML reads config information from the YAML file at the specified path
func (c *Config) ReadFromYAML(fp string) error {
	configfileBytes, err := ioutil.ReadFile(fp)
	yaml.Unmarshal(configfileBytes, c)
	c.BinDir = path.Join(c.WebpageDir, c.BinPath)
	return err
}

//ReadFromEnvVars reads config information from environment variables
//Whatever you do, never stop using this. This took WAY too long to write.
func (c *Config) ReadFromEnvVars() error {
	cType := reflect.ValueOf(*c).Type()
	cElem := reflect.ValueOf(c).Elem()

	for i:=0; i<reflect.ValueOf(*c).NumField(); i++ {
		field := cType.Field(i)
		fieldVal := cElem.Field(i)
		tag, tagExists := field.Tag.Lookup("env")
		if tagExists {
			envVar, envVarExists := os.LookupEnv(tag)
			if envVarExists {
				var err error
				switch field.Type.Kind() {
				case reflect.Int:
					val, e := strconv.Atoi(envVar)
					fieldVal.Set(reflect.ValueOf(val))
					err = e
				case reflect.Int64:
					val, e := strconv.ParseInt(envVar, 10, 64) // strconv.Atoi can't return int64
					fieldVal.Set(reflect.ValueOf(val))
					err = e
				case reflect.String:
					fieldVal.Set(reflect.ValueOf(envVar))
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
