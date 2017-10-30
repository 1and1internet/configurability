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

type NginxJsonData struct {
	Gzip	string	`json:"gzip"`
	DocumentRoot	string	`json:"document_root"`
}

type NginxParserData struct {
	JsonData NginxJsonData
	Section ini.Section
	TestInputFolder string
	TestOutputFolder string
	NginxSourceDirectory string
	NginxDestDirectory string
	GzipSourceFilePath string
	GzipDestFilePath string
	GzipLevel string
	GzipState string
}

func (nginx *NginxParserData) NginxJsonLoader(data []byte) {
	err := json.Unmarshal(data, &nginx.JsonData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	nginx.TestInputFolder = os.Getenv("TEST_INPUT_FOLDER")
	nginx.TestOutputFolder = os.Getenv("TEST_OUTPUT_FOLDER")
	nginx.NginxSourceDirectory =  path.Join(nginx.TestInputFolder, "/etc", "nginx")
	nginx.NginxDestDirectory = path.Join(nginx.TestOutputFolder, "/etc", "nginx")
	nginx.GzipSourceFilePath = path.Join(nginx.NginxSourceDirectory, "conf.d", "gzip.conf")
	nginx.GzipDestFilePath = path.Join(nginx.NginxDestDirectory, "conf.d", "gzip.conf")
	if strings.ToLower(nginx.JsonData.Gzip) == "off" {
		nginx.GzipLevel = "6"
		nginx.GzipState = "off"
	} else {
		nginx.GzipLevel = nginx.JsonData.Gzip
		nginx.GzipState = "on"
	}
}

func (nginx *NginxParserData) DocrootFix()  {
	if nginx.JsonData.DocumentRoot != "" {
		sites_enabled_source_directory := path.Join(nginx.NginxSourceDirectory, "sites-enabled")
		sites_enabled_dest_directory := path.Join(nginx.NginxDestDirectory, "sites-enabled")

		document_root_key := "DOCUMENT_ROOT"
		document_root_default := "html"
		document_root := nginx.JsonData.DocumentRoot

		envar_current, ok := os.LookupEnv(document_root_key)
		if ok && envar_current != document_root && envar_current != document_root_default {
			log.Printf("'Legacy %s variable is present with a conflicting value", document_root_key)
			return
		}

		// Create the document root folder if it's missing
		document_root_path := path.Join(nginx.TestOutputFolder, "/var", "www", document_root)
		EnsureDirExists(document_root_path)

		variable_regex := regexp.MustCompile(fmt.Sprintf("\\${?%s}?", document_root_key))
		root_command_regex := regexp.MustCompile("root /var/www/.*;")
		new_root_command := fmt.Sprintf("root /var/www/%s;", document_root)

		files, err := ioutil.ReadDir(sites_enabled_source_directory)
		if err != nil {
			log.Printf("The following error occurred while trying to list the sites-enabled folder: %s", err)
			return
		}

		for _, file_path := range files {
			write_needed := false
			new_lines := []string{}
			full_file_path := path.Join(sites_enabled_source_directory, file_path.Name())
			current_lines := ReadLinesFromFile(full_file_path)

			for _, line := range current_lines {
				new_line := variable_regex.ReplaceAllString(line, document_root)
				new_line = root_command_regex.ReplaceAllString(new_line, new_root_command)
				new_lines = append(new_lines, new_line)
				write_needed = write_needed || line != new_line
			}

			if write_needed {
				full_file_path := path.Join(sites_enabled_dest_directory, file_path.Name())
				WriteLinesToFile(full_file_path, new_lines)
			}
		}
	}
}

func (nginx *NginxParserData) GzipFix() {
	if nginx.JsonData.Gzip != "" {
		new_gzip_command := fmt.Sprintf("gzip %s;", nginx.GzipState)
		new_gzip_level_command := fmt.Sprintf("gzip_comp_level %s;", nginx.GzipLevel)

		gzip_command_regex := regexp.MustCompile("gzip \\w*;")
		gzip_level_command_regex := regexp.MustCompile("gzip_comp_level \\d*;")

		current_lines := ReadLinesFromFile(nginx.GzipSourceFilePath)
		new_lines := []string{}
		write_needed := false

		for _, line := range current_lines {
			new_line := gzip_command_regex.ReplaceAllString(line, new_gzip_command)
			new_line = gzip_level_command_regex.ReplaceAllString(new_line, new_gzip_level_command)
			new_lines = append(new_lines, new_line)
			write_needed = write_needed || line != new_line
		}

		if write_needed {
			WriteLinesToFile(nginx.GzipDestFilePath, new_lines)
		}
	}
}

func (nginx *NginxParserData) ApplyCustomisations() {
	enabled, _, _ := unpack_etc_ini(nginx.Section, false)
	if enabled {
		nginx.DocrootFix()
		nginx.GzipFix()
	}
}

