package main

import (
	"C"
	"io/ioutil"
	"log"

	"github.com/1and1internet/configurability/plugins"
	"github.com/go-ini/ini"
	"gopkg.in/yaml.v2"
)
import (
	"os"
	"path"
	"strings"
)

type Customisor interface {
	ApplyCustomisations()
}

type MongoSystemLogConfig struct {
	Destination string `yaml:"destination"`
	LogAppend   bool   `yaml:"logAppend"`
	Path        string `yaml:"path"`
}

type MongoStorageJournalConfig struct {
	Enabled bool `yaml:"enabled"`
}

type MongoStorageConfig struct {
	DbPath  string                    `yaml:"dbPath"`
	Journal MongoStorageJournalConfig `yaml:"journal"`
}

type MongoNetConfig struct {
	Port   int    `yaml:"port"`
	BindIp string `yaml:"bindIp"`
}

type MongoProcessManagementConfig struct {
	TimeZoneInfo string `yaml:"timeZoneInfo"`
}

type MongoSecurityConfig struct {
}

type MongoOperationProfilingConfig struct {
}

type MongoYamlData struct {
	Storage            MongoStorageConfig            `yaml:"storage"`
	SystemLog          MongoSystemLogConfig          `yaml:"systemLog"`
	Net                MongoNetConfig                `yaml:"net"`
	ProcessManagement  MongoProcessManagementConfig  `yaml:"processManagement"`
	Security           MongoSecurityConfig           `yaml:"security"`
	OperationProfiling MongoOperationProfilingConfig `yaml:"operationProfiling"`
}

type MongoData struct {
	Config                *MongoYamlData
	ConfMongodJsonSection ini.Section
	SourceConfigFilePath  string
}

func (mongoConfig *MongoData) LoadCustomConfig(data []byte) {
	err := yaml.Unmarshal(data, mongoConfig.Config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func (mongoData *MongoData) LoadConfig(section ini.Section) bool {
	var sourceConfigFilePath = ""
	if plugins.GetFromSection(section, "enabled", "false", true) != "true" {
		return false
	}
	sourceConfigFilePath = plugins.GetFromSection(section, "ini_file_path", "", false)
	if sourceConfigFilePath == "" {
		return false
	}
	mongoData.SourceConfigFilePath = sourceConfigFilePath

	iniFileData, err := ioutil.ReadFile(sourceConfigFilePath)
	if err != nil {
		log.Printf("Error loading default config: %v ", err)
		return false
	}

	err = yaml.Unmarshal(iniFileData, mongoData.Config)
	if err != nil {
		log.Fatalf("Error unmarshalling default config: %v", err)
		return false
	}

	return true
}

func (data *MongoData) Save() {
	var test_output_folder = os.Getenv("TEST_OUTPUT_FOLDER")
	var target_output_file = data.SourceConfigFilePath
	if test_output_folder != "" {
		target_output_file = strings.Join([]string{test_output_folder, path.Base(data.SourceConfigFilePath)}, "/")
	}

	yamlData, err := yaml.Marshal(data.Config)
	if err != nil {
		log.Printf("Error marshalling data for yaml: %v", err)
		return
	}
	log.Printf("Writing mongodb config to %s:  \n\n%s\n", target_output_file, string(yamlData))
	yamlDataList := strings.Split(string(yamlData), "\n")
	plugins.WriteLinesToFile(target_output_file, yamlDataList)
}

func (data *MongoData) ApplyCustomisations(content []byte) {
	enabled := data.LoadConfig(data.ConfMongodJsonSection)
	if enabled {
		data.LoadCustomConfig(content)
		data.Save()
	}
}

func Customise(content []byte, section *ini.Section, configurationFileName string) bool {
	if configurationFileName == "configuration-mongod.json" {
		log.Println("Process as mongo/yaml")
		data := MongoData{
			Config:                &MongoYamlData{},
			ConfMongodJsonSection: *section,
		}
		data.ApplyCustomisations(content)
		return true
	}
	return false
}
