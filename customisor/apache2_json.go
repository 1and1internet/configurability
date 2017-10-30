package customisor

import (
	"github.com/go-ini/ini"
	"encoding/json"
	"log"
	"path"
	"os"
	"regexp"
	"fmt"
	"io/ioutil"
	"strings"
)

type ApacheJsonData struct {
	Gzip	string	`json:"gzip"`
	DocumentRoot	string	`json:"document_root"`
}

type ApacheParserData struct {
	JsonData ApacheJsonData
	Section ini.Section
	TestInputFolder string
	TestOutputFolder string
	ApacheSourceDirectory string
	ApacheDestDirectory string
	GzipSourceFilePath string
	GzipDestFilePath string
	GzipLevel string
	GzipState string
	DocumentRootKey string
}

func (apache *ApacheParserData) ApacheJsonLoader(data []byte) {
	err := json.Unmarshal(data, &apache.JsonData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	apache.TestInputFolder = os.Getenv("TEST_INPUT_FOLDER")
	apache.TestOutputFolder = os.Getenv("TEST_OUTPUT_FOLDER")
	apache.ApacheSourceDirectory =  path.Join(apache.TestInputFolder, "/etc", "apache2")
	apache.ApacheDestDirectory = path.Join(apache.TestOutputFolder, "/etc", "apache2")
	apache.GzipSourceFilePath = path.Join(apache.ApacheSourceDirectory, "conf.d", "gzip.conf")
	apache.GzipDestFilePath = path.Join(apache.ApacheDestDirectory, "conf.d", "gzip.conf")
	if strings.ToLower(apache.JsonData.Gzip) == "off" {
		apache.GzipLevel = "6"
		apache.GzipState = "off"
	} else {
		apache.GzipLevel = apache.JsonData.Gzip
		apache.GzipState = "on"
	}
	apache.DocumentRootKey = "DOCUMENT_ROOT"
}

func (apache *ApacheParserData) updateApacheConfig(source, dest string) {
	regex := regexp.MustCompile(fmt.Sprintf("\\${?%s}?", apache.DocumentRootKey))
	files, err := ioutil.ReadDir(source)
	if err != nil {
		log.Printf("The following error occurred while trying to list the sites-enabled folder: %s", err)
		return
	}

	for _, file_path := range files {
		write_needed := false
		new_lines := []string{}
		full_file_path := path.Join(source, file_path.Name())
		current_lines := ReadLinesFromFile(full_file_path)

		for _, line := range current_lines {
			new_line := regex.ReplaceAllString(line, apache.JsonData.DocumentRoot)
			new_lines = append(new_lines, new_line)
			write_needed = write_needed || line != new_line
		}

		if write_needed {
			full_file_path := path.Join(dest, file_path.Name())
			WriteLinesToFile(full_file_path, new_lines)
		}
	}
}

func (apache *ApacheParserData) DocrootFix()  {
	if apache.JsonData.DocumentRoot != "" {
		sites_enabled_source_directory := path.Join(apache.ApacheSourceDirectory, "sites-enabled")
		sites_enabled_dest_directory := path.Join(apache.ApacheDestDirectory, "sites-enabled")
		conf_enabled_source_directory := path.Join(apache.ApacheSourceDirectory, "conf-enabled")
		conf_enabled_dest_directory := path.Join(apache.ApacheDestDirectory, "conf-enabled")
		mods_enabled_source_directory := path.Join(apache.ApacheSourceDirectory, "mods-enabled")
		mods_enabled_dest_directory := path.Join(apache.ApacheDestDirectory, "mods-enabled")

		document_root_default := "html"
		document_root := apache.JsonData.DocumentRoot

		envar_current, ok := os.LookupEnv(apache.DocumentRootKey)
		if ok && envar_current != document_root && envar_current != document_root_default {
			log.Printf("'Legacy %s variable is present with a conflicting value", apache.DocumentRootKey)
			return
		}

		// Create the document root folder if it's missing
		document_root_path := path.Join(apache.TestOutputFolder, "/var", "www", document_root)
		EnsureDirExists(document_root_path)

		EnsureDirExists(sites_enabled_dest_directory)
		EnsureDirExists(conf_enabled_dest_directory)
		EnsureDirExists(mods_enabled_dest_directory)

		apache.updateApacheConfig(sites_enabled_source_directory, sites_enabled_dest_directory)
		apache.updateApacheConfig(conf_enabled_source_directory, conf_enabled_dest_directory)
		apache.updateApacheConfig(mods_enabled_source_directory, mods_enabled_dest_directory)
	}
}

func (apache *ApacheParserData) GzipFix() {
	if apache.JsonData.Gzip != "" {
		mods_enabled_dest_directory := path.Join(apache.ApacheDestDirectory, "mods-enabled")
		EnsureDirExists(mods_enabled_dest_directory)
		if apache.GzipState == "off" {
			for _, file_path := range []string{"deflate.conf", "deflate.load"} {
				full_file_path := path.Join(mods_enabled_dest_directory, file_path)
				_, err := os.Stat(full_file_path)
				if err == nil {
					remove_err := os.Remove(full_file_path)
					if remove_err != nil {
						log.Printf("Could not remove %s: %s", full_file_path, remove_err)
					}
				}
			}
		} else {
			full_file_path := path.Join(mods_enabled_dest_directory, "deflate.conf")
			fh, err := os.Create(full_file_path)
			if err != nil {
				log.Printf("Could not write to %s: %s", full_file_path, err)
				return
			}
			fh.WriteString(fmt.Sprintf("DeflateCompressionLevel %s\n", apache.GzipLevel))
		}
	}
}

func (apache *ApacheParserData) ApplyCustomisations() {
	enabled, _, _ := unpack_etc_ini(apache.Section, false)
	if enabled {
		apache.DocrootFix()
		apache.GzipFix()
	}
}
