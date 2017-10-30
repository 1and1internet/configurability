package file_helpers

import (
	"os"
	"log"
	"io/ioutil"
	"configurability/customisor"
	"path"
)

type CustomisationData struct {
	Basic  customisor.BasicYamlData
	MyDb   customisor.MysqlParserData
	Php    customisor.PhpParserData
	Nginx  customisor.NginxParserData
	Apache customisor.ApacheParserData
}

func get_customisation_folder() (string) {
	var confDir = os.Getenv("CONFIGURABILITY_DIR")
	return confDir
}

func get_folder_content_list(folder string) ([]string) {
	var filePaths []string
	if folder != "" {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			filePaths = append(filePaths, path.Join(folder, file.Name()))
		}
	}
	return filePaths
}

func get_folder_content_map(folder string) (map[string]string) {
	fileMap := make(map[string]string)
	if folder != "" {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fileMap[file.Name()] = path.Join(folder, file.Name())
		}
	}
	return fileMap
}

func get_etc_config_folder() (string) {
	var confDir = os.Getenv("CONFIGURABILITY_INTERNAL")
	if confDir == "" {
		confDir = "/etc/configurability/"
	}
	return confDir
}

func list_etc_config_folder() ([]string) {
	return get_folder_content_list(get_etc_config_folder())
}

func map_customisation_folder() (map[string]string) {
	return get_folder_content_map(get_customisation_folder())
}

func (cd *CustomisationData) LoadCustomisationData() {
	var customisationFilePathMap map[string]string = map_customisation_folder()
	for _, etcConfigrationFilePath := range list_etc_config_folder() {
		var section = customisor.ReadEtcConfiguration(etcConfigrationFilePath)
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

				if configuration_file_name.String() == "configuration-basic.yaml" {
					log.Println("Process as basic/yaml")
					var basic customisor.BasicYamlData
					basic.BasicYamlLoader(content)
				} else if configuration_file_name.String() == "configuration-mysql.json" {
					log.Println("Process as mysql/json")
					cd.MyDb = customisor.MysqlParserData{}
					cd.MyDb.MysqlJsonLoader(content)
					cd.MyDb.Section = *section
					cd.MyDb.ApplyCustomisations()
				} else if configuration_file_name.String() == "configuration-nginx.json" {
					log.Println("Process as nginx/json")
					cd.Nginx = customisor.NginxParserData{}
					cd.Nginx.NginxJsonLoader(content)
					cd.Nginx.Section = *section
					cd.Nginx.ApplyCustomisations()
				} else if configuration_file_name.String() == "configuration-apache2.json" {
					log.Println("Process as apache2/json")
					cd.Apache = customisor.ApacheParserData{}
					cd.Apache.ApacheJsonLoader(content)
					cd.Apache.Section = *section
					cd.Apache.ApplyCustomisations()
				} else if configuration_file_name.String() == "configuration-php.json" {
					log.Println("Process as php/json")
					cd.Php = customisor.PhpParserData{}
					cd.Php.PhpJsonLoader(content)
					cd.Php.Section = *section
					cd.Php.ApplyCustomisations()
				} else {
					log.Printf("WARNING: Unexpected filename", customisationFilePath)
				}
			}
		}
	}
}