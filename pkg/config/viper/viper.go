// Copyright © 2019 Pavel Hadzhiev <p.hadzhiev96@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package viper

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/pavelhadzhiev/story-builder/pkg/config"
)

// DefaultConfigFileName is the default configuration path, in case no other has been provided.
const DefaultConfigFileName = ".story-builder.json"

// Configurator implements the SBConfigurator interface using viper and the file system.
type Configurator struct {
	viper *viper.Viper
}

// NewConfigurator creates a new ViperConfigurator, provided a configuration file path.
func NewConfigurator(cfgFile string) (config.SBConfigurator, error) {
	viper := viper.New()

	absCfgFilePath, err := getConfigFileAbsPath(cfgFile)
	if err != nil {
		return nil, err
	}
	viper.SetConfigFile(absCfgFilePath)

	configurator := &Configurator{viper: viper}
	if err := viper.ReadInConfig(); err != nil {
		configurator.Save(&config.SBConfiguration{})
	}

	return configurator, nil
}

// Save is filling the configuration file with the properties from the passed configuration object
func (viperConfig *Configurator) Save(sbConfig *config.SBConfiguration) error {
	viperConfig.viper.Set("url", sbConfig.URL)
	viperConfig.viper.Set("authorization", sbConfig.Authorization)
	viperConfig.viper.Set("room", sbConfig.Room)

	if err := viperConfig.viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

// Load returns a SBConfiguration pointer, storing the properties from the configuration file
func (viperConfig *Configurator) Load() (*config.SBConfiguration, error) {
	if err := viperConfig.viper.ReadInConfig(); err != nil {
		return nil, err
	}

	sbConfig := &config.SBConfiguration{}
	if err := viperConfig.viper.Unmarshal(sbConfig); err != nil {
		return nil, err
	}

	return sbConfig, nil
}

func getConfigFileAbsPath(cfgFile string) (string, error) {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		cfgFile = filepath.Join(home, DefaultConfigFileName)
	}

	filename, err := filepath.Abs(cfgFile)
	if err != nil {
		return "", err
	}

	return filename, nil
}
