package main

import (
	"C"
	"github.com/go-ini/ini"
	"encoding/json"
	"log"
	"github.com/1and1internet/configurability/plugins"
)

type PhpJsonData struct {
	PHP struct {
		Engine                bool   `json:"engine"`
		FileUploads           bool   `json:"file_uploads"`
		ShortOpenTag          bool   `json:"short_open_tag"`
		EnablePostDataReading bool   `json:"enable_post_data_reading"`
		AllowUrlFopen         bool   `json:"allow_url_fopen"`
		AllowUrlInclude       bool   `json:"allow_url_include"`
		OutputBuffering       string `json:"output_buffering"`
		MaxInputTime          int64  `json:"max_input_time"`
		MaxExecutionTime      int64  `json:"max_execution_time"`
		MaxInputVars          string `json:"max_input_vars"`
		MaxInputNestingLevel  string `json:"max_input_nesting_level"`
		MemoryLimit           string `json:"memory_limit"`
		PostMaxSize           string `json:"post_max_size"`
		UploadMaxFilesize     string `json:"upload_max_filesize"`
		MaxFileUploads        string `json:"max_file_uploads"`
		DisplayErrors         bool   `json:"display_errors"`
		DisplayStartupErrors  bool   `json:"display_startup_errors"`
		HtmlErrors            bool   `json:"html_errors"`
		LogErrors             bool   `json:"log_errors"`
		IgnoreRepeatedErrors  bool   `json:"ignore_repeated_errors"`
		TrackErrors           bool   `json:"track_errors"`
		ErrorReporting        string `json:"error_reporting"`
	}
	OpCache struct {
		OpCacheMemory         string `json:"php_opcache.memory_consumption"`
	}
}

type PhpParserData struct {
	JsonData PhpJsonData
	Section ini.Section
}

func (php *PhpParserData) PhpJsonLoader(data []byte) {
	err := json.Unmarshal(data, &php.JsonData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func (php *PhpParserData) ApplyCustomisations() {
	_, iniFile, iniFilePath := plugins.UnpackEtcIni(php.Section, false)
	if iniFile != nil {
		php_section, err := iniFile.GetSection("PHP")
		if err == nil {
			plugins.UpdateBoolKey("PHP", php_section, "engine", php.JsonData.PHP.Engine)
			plugins.UpdateBoolKey("PHP", php_section, "file_uploads", php.JsonData.PHP.FileUploads)
			plugins.UpdateBoolKey("PHP", php_section, "short_open_tag", php.JsonData.PHP.ShortOpenTag)
			plugins.UpdateBoolKey("PHP", php_section, "enable_post_data_reading", php.JsonData.PHP.EnablePostDataReading)
			plugins.UpdateBoolKey("PHP", php_section, "allow_url_fopen", php.JsonData.PHP.AllowUrlFopen)
			plugins.UpdateBoolKey("PHP", php_section, "allow_url_include", php.JsonData.PHP.AllowUrlInclude)

			plugins.UpdateStringKey("PHP", php_section, "output_buffering", php.JsonData.PHP.OutputBuffering)

			plugins.UpdateInt64Key("PHP", php_section, "max_input_time", php.JsonData.PHP.MaxInputTime)
			plugins.UpdateInt64Key("PHP", php_section, "max_execution_time", php.JsonData.PHP.MaxExecutionTime)

			plugins.UpdateStringKey("PHP", php_section, "max_input_vars", php.JsonData.PHP.MaxInputVars)
			plugins.UpdateStringKey("PHP", php_section, "max_input_nesting_level", php.JsonData.PHP.MaxInputNestingLevel)
			plugins.UpdateStringKey("PHP", php_section, "memory_limit", php.JsonData.PHP.MemoryLimit)
			plugins.UpdateStringKey("PHP", php_section, "post_max_size", php.JsonData.PHP.PostMaxSize)
			plugins.UpdateStringKey("PHP", php_section, "upload_max_filesize", php.JsonData.PHP.UploadMaxFilesize)
			plugins.UpdateStringKey("PHP", php_section, "max_file_uploads", php.JsonData.PHP.MaxFileUploads)

			plugins.UpdateBoolKey("PHP", php_section, "display_errors", php.JsonData.PHP.DisplayErrors)
			plugins.UpdateBoolKey("PHP", php_section, "display_startup_errors", php.JsonData.PHP.DisplayStartupErrors)
			plugins.UpdateBoolKey("PHP", php_section, "html_errors", php.JsonData.PHP.HtmlErrors)
			plugins.UpdateBoolKey("PHP", php_section, "log_errors", php.JsonData.PHP.LogErrors)
			plugins.UpdateBoolKey("PHP", php_section, "ignore_repeated_errors", php.JsonData.PHP.IgnoreRepeatedErrors)
			plugins.UpdateBoolKey("PHP", php_section, "track_errors", php.JsonData.PHP.TrackErrors)

			plugins.UpdateStringKey("PHP", php_section, "error_reporting", php.JsonData.PHP.ErrorReporting)
		}
		
		opcache_section, err := iniFile.GetSection("php_opcache")
		if err == nil {
			plugins.UpdateStringKey("PHP", opcache_section, "php_opcache.memory_consumption", php.JsonData.OpCache.OpCacheMemory)
		}

			plugins.SaveIniFile(*iniFile, iniFilePath, "php.ini")
	}
}

func Customise(content []byte, section *ini.Section, configurationFileName string) (bool) {
	if configurationFileName == "configuration-php.json" {
		log.Println("Process as php/json")
		php := PhpParserData{}
		php.PhpJsonLoader(content)
		php.Section = *section
		php.ApplyCustomisations()
		return true
	}
	return false
}
