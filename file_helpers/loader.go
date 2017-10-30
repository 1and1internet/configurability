package file_helpers

import (
	"os"
	"log"
	"io/ioutil"
	"path"
)

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

func ListEtcConfigFolder() ([]string) {
	return get_folder_content_list(get_etc_config_folder())
}

func MapCustomisationFolder() (map[string]string) {
	return get_folder_content_map(get_customisation_folder())
}