package cmd

import (
	"os"

	"fmt"

	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func SaveViperConfig() error {
	cfg := viper.ConfigFileUsed()
	if cfg == "" {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		cfg = filepath.Join(usr.HomeDir, "."+NAME+".yaml")
	}
	f, err := os.Create(cfg)
	if err != nil {
		return err
	}
	defer f.Close()

	all := viper.AllSettings()
	b, err := yaml.Marshal(all)
	if err != nil {
		return fmt.Errorf("Panic while encoding into YAML format.")
	}
	if _, err := f.WriteString(string(b)); err != nil {
		return err
	}
	return nil
}
