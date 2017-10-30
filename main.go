package main

import (
	"path/filepath"
	"path"
	"log"
	"plugin"
	"os"
	"github.com/1and1internet/configurability/file_helpers"
	"github.com/1and1internet/configurability/plugins"
	"github.com/go-ini/ini"
	"io/ioutil"
)

type CustomisorInterface interface {
	Customise([]byte, *ini.Section, string) bool
}

func getPluginFolder() (string) {
	pluginFolder, ok := os.LookupEnv("CONF_PLUGIN_FOLDER")
	if ok {
		return pluginFolder
	}
	return "/opt/configurability/goplugins"
}

func LoadCustomisationData(customisorSymbol plugin.Symbol) {
	var customisationFilePathMap map[string]string = file_helpers.MapCustomisationFolder()
	found := false
	for _, etcConfigrationFilePath := range file_helpers.ListEtcConfigFolder() {
		var section = plugins.ReadEtcConfiguration(etcConfigrationFilePath)
		if section != nil {
			var configuration_file_name = section.Key("configuration_file_name")
			customisationFilePath, ok := customisationFilePathMap[configuration_file_name.String()]
			if ok {
				content, err := ioutil.ReadFile(customisationFilePath)
				if err != nil {
					log.Printf("There was a problem reading %s: %s\n", configuration_file_name.String(), err)
					log.Println("Continuing without it...")
					continue
				}

				found = customisorSymbol.(func([]byte, *ini.Section, string) bool)(content, section, configuration_file_name.String())
				if found {
					break
				}
			}
		}
	}
	if !found {
		log.Printf("WARNING: No plugin for symbol %v", customisorSymbol)
	}
}

func main() {
	loggingFilename := "/dev/stdout"
	f, err := os.OpenFile(loggingFilename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
			log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	pluginFolder := getPluginFolder()
	fileglob := path.Join(pluginFolder, "*.so")
	files, err := filepath.Glob(fileglob)
	if err == nil {
		for _, file := range files {
			log.Printf("Loading plugin %s\n", file)
			configuratorPlugin, err := plugin.Open(file)
			if err != nil {
				log.Printf("Could not load plugin %s: %s\n", file, err)
				continue
			}
			customisorSymbol, err := configuratorPlugin.Lookup("Customise")
			if err != nil {
				log.Printf("Could not lookup 'Customise' in %s\n", file)
			}
			LoadCustomisationData(customisorSymbol)
		}
	} else {
		log.Printf("Fileglob error: %s\n", err)
	}
}
