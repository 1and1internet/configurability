package main

import (
	"C"
	"github.com/go-ini/ini"
	"encoding/json"
	"log"
	"github.com/1and1internet/configurability/plugins"
)

type OpcacheJsonData struct {
	OpCache struct {
		OpCacheMemory         string `json:"php_opcache.memory_consumption"`
		OpCacheRevalidateFreq int64 `json:"php_opcache.revalidate_freq"`
		OpCacheEnableCli	  bool `json:"php_opcache.enable_cli"`
	}
}

type OpcacheParserData struct {
	JsonData OpcacheJsonData
	Section  ini.Section
}

func (opcache *OpcacheParserData) OpcacheJsonLoader(data []byte) {
	// Set some defaults
	opcache.JsonData.OpCache.OpCacheMemory = "128"
	opcache.JsonData.OpCache.OpCacheRevalidateFreq = 2
	opcache.JsonData.OpCache.OpCacheEnableCli = false
	err := json.Unmarshal(data, &opcache.JsonData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func (opcache *OpcacheParserData) ApplyCustomisations() {
	_, iniFile, iniFilePath := plugins.UnpackEtcIni(opcache.Section, false)
	if iniFile != nil {
		opcache_section, err := iniFile.GetSection("")
		if err == nil {
			plugins.UpdateStringKey("Opcache", opcache_section, "opcache.memory_consumption", opcache.JsonData.OpCache.OpCacheMemory)
			plugins.UpdateInt64Key("Opcache", opcache_section, "opcache.revalidate_freq", opcache.JsonData.OpCache.OpCacheRevalidateFreq)
			plugins.UpdateBoolKey("Opcache", opcache_section, "opcache.enable_cli", opcache.JsonData.OpCache.OpCacheEnableCli)
		}
		plugins.SaveIniFile(*iniFile, iniFilePath, "10-opcache.ini")
	}
}

func Customise(content []byte, section *ini.Section, configurationFileName string) (bool) {
	if configurationFileName == "configuration-php.json" {
		log.Println("Process as php-opcache/json")
		php := OpcacheParserData{}
		php.OpcacheJsonLoader(content)
		php.Section = *section
		php.ApplyCustomisations()
		return true
	}
	return false
}
