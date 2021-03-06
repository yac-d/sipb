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
//Whatever you do, never stop using this, because it took WAY too long to write.
func (c *Config) ReadFromEnvVars() (err error) {
	cVal := reflect.ValueOf(*c)
	cType := cVal.Type()
	cElem := reflect.ValueOf(c).Elem()

	for i:=0; i<cVal.NumField(); i++ {
		field := cType.Field(i)
		fieldVal := cElem.Field(i)
		tag, tagExists := field.Tag.Lookup("env")
		if tagExists {
			envVar, envVarExists := os.LookupEnv(tag)
			if envVarExists {
				var val any
				switch field.Type.Kind() {
				case reflect.Int:
					val, err = strconv.Atoi(envVar)
				case reflect.Int64:
					val, err = strconv.ParseInt(envVar, 10, 64) // strconv.Atoi can't return int64
				case reflect.String:
					val = envVar
				}
				fieldVal.Set(reflect.ValueOf(val))
			}
		}
	}
	c.BinDir = path.Join(c.WebpageDir, c.BinPath)

	return
}
