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
	"strconv"
	"strings"
)

type Customisor interface {
	ApplyCustomisations()
}

type MongoSystemLogConfig struct {
	Destination        string                        `yaml:"destination"`
	LogAppend          bool                          `yaml:"logAppend"`
	Path               string                        `yaml:"path"`
	Verbosity          int                           `yaml:"verbosity"`
	Quiet              bool                          `yaml:"quiet"`
	TraceAllExceptions bool                          `yaml:"traceAllExceptions"`
	SyslogFacility     string                        `yaml:"syslogFacility"`
	LogRotate          string                        `yaml:"logRotate"`
	TimeStampFormat    string                        `yaml:"timeStampFormat"`
	Component          MongoSystemLogComponentConfig `yaml:"component"`
}

type MongoSystemLogComponentConfig struct {
	AccessControl MongoSystemLogComponentVerbosityConfig `yaml:"accessControl"`
	Command       MongoSystemLogComponentVerbosityConfig `yaml:"command"`
	// Replication not defined because we wont need it (yet)
	// Storage not defined because we wont need it (yet)
	// Write not defined because we wont need it (yet)
}

type MongoSystemLogComponentVerbosityConfig struct {
	Verbosity int `yaml:"verbosity"`
}

type MongoStorageConfig struct {
	DbPath          string                       `yaml:"dbPath"`
	Journal         MongoStorageJournalConfig    `yaml:"journal"`
	IndexBuildRetry bool                         `yaml:"indexBuildRetry"`
	WiredTiger      MongoStorageWiredTigerConfig `yaml:"wiredTiger"`
}

type MongoStorageJournalConfig struct {
	Enabled bool `yaml:"enabled"`
}

type MongoStorageWiredTigerConfig struct {
	EngineConfig MongoStorageWiredTigerEngingConfig `yaml:"engineConfig"`
}

type MongoStorageWiredTigerEngingConfig struct {
	CacheSizeGB float32 `yaml:"cacheSizeGB"`
}

type MongoNetConfig struct {
	Port            int                 `yaml:"port"`
	BindIpAll       bool                `yaml:"bindIpAll"`
	WireObjectCheck bool                `yaml:"wireObjectCheck"`
	Ipv6            bool                `yaml:"ipv6"`
	Ssl             MongoNetSslConfig   `yaml:"ssl"`
	Compression     MongoNetCompression `yaml:"compression"`
	ServiceExecutor string              `yaml:"serviceExecutor"`
}

type MongoNetAltConfig struct {
	Port            int                 `yaml:"port"`
	BindIpAll       bool                `yaml:"bindIpAll"`
	WireObjectCheck bool                `yaml:"wireObjectCheck"`
	Ipv6            bool                `yaml:"ipv6"`
	Ssl             MongoNetNoSslConfig `yaml:"ssl"`
	Compression     MongoNetCompression `yaml:"compression"`
	ServiceExecutor string              `yaml:"serviceExecutor"`
}

type MongoNetSslConfig struct {
	Mode                                string `yaml:"mode"`
	PEMKeyFile                          string `yaml:"PEMKeyFile"`
	PEMKeyPassword                      string `yaml:"PEMKeyPassword"`
	ClusterFile                         string `yaml:"clusterFile"`
	ClusterPassword                     string `yaml:"clusterPassword"`
	CAFile                              string `yaml:"CAFile"`
	CRLFile                             string `yaml:"CRLFile"`
	AllowConnectionsWithoutCertificates bool   `yaml:"allowConnectionsWithoutCertificates"`
	AllowInvalidCertificates            bool   `yaml:"allowInvalidCertificates"`
	AllowInvalidHostnames               bool   `yaml:"allowInvalidHostnames"`
	DisabledProtocols                   string `yaml:"disabledProtocols"`
	FIPSMode                            bool   `yaml:"FIPSMode"`
}

type MongoNetNoSslConfig struct {
	Mode string `yaml:"mode"`
}

type MongoNetCompression struct {
	Compressors string `yaml:"compressors"`
}

type MongoProcessManagementConfig struct {
	TimeZoneInfo string `yaml:"timeZoneInfo"`
	Fork         bool   `yaml:"fork"`
	PidFilePath  string `yaml:"pidFilePath"`
}

type MongoSecurityConfig struct {
	Authorization     string `yaml:"authorization"`
	JavascriptEnabled bool   `yaml:"javascriptEnabled"`
}

type MongoYamlData struct {
	Storage MongoStorageConfig `yaml:"storage"`
	// SystemLog          MongoSystemLogConfig          `yaml:"systemLog"`
	Net MongoNetConfig `yaml:"net"`
	// ProcessManagement  MongoProcessManagementConfig  `yaml:"processManagement"`
	Security MongoSecurityConfig `yaml:"security"`
}

type MongoYamlNoSslData struct {
	Storage  MongoStorageConfig  `yaml:"storage"`
	Net      MongoNetAltConfig   `yaml:"net"`
	Security MongoSecurityConfig `yaml:"security"`
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

func (mongoData *MongoData) ApplyLocalStorageWiredTigerConfig() {
	mongoData.Config.Storage.WiredTiger.EngineConfig.CacheSizeGB = 0.25
	cgroup_mem_limit_str := plugins.GetMaxMemoryOfContainerAsString()
	cgroup_mem_limit_int, err := strconv.Atoi(cgroup_mem_limit_str)
	if err == nil {
		new_limit := float32((cgroup_mem_limit_int / 1024 / 1024 / 2) - 1024)
		if new_limit > mongoData.Config.Storage.WiredTiger.EngineConfig.CacheSizeGB {
			mongoData.Config.Storage.WiredTiger.EngineConfig.CacheSizeGB = new_limit
		}
	}
}

func (mongoData *MongoData) ApplyLocalStorageConfig() {
	mongoData.ApplyLocalStorageWiredTigerConfig()
}

func (mongoData *MongoData) ApplyLocalEnvironmentConfig() {
	mongoData.ApplyLocalStorageConfig()
}

func (data *MongoData) NoSslConversion(yamlData []byte) []byte {
	// We don't even use SSL yet, but if we ever do then this will be useful.
	// When ssl is enabled, the rest of the ssl parameters are used. If they are
	// there when ssl is disabled then mongo complains, so we repackage without them.
	if data.Config.Net.Ssl.Mode == "disabled" {
		noSslData := &MongoYamlNoSslData{}
		err := yaml.Unmarshal(yamlData, noSslData)
		if err != nil {
			log.Fatalf("nossl error: %v", err)
		}
		yamlData, err = yaml.Marshal(noSslData)
		if err != nil {
			log.Printf("Error marshalling data for yaml (nossl): %v", err)
		}
	}
	return yamlData
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
	yamlData = data.NoSslConversion(yamlData)

	log.Printf("Writing mongodb config to %s:  \n\n%s\n", target_output_file, string(yamlData))
	yamlDataList := strings.Split(string(yamlData), "\n")
	plugins.WriteLinesToFile(target_output_file, yamlDataList)
}

func (data *MongoData) ApplyCustomisations(content []byte) {
	enabled := data.LoadConfig(data.ConfMongodJsonSection)
	if enabled {
		data.ApplyLocalEnvironmentConfig()
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
