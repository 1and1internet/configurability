package plugins

import (
	"github.com/go-ini/ini"
	"strings"
	"log"
	"os"
	"strconv"
	"bufio"
	"path"
	"fmt"
)

func get_from_section(iniSection ini.Section, key string, defaultValue string, doToLower bool) (string) {
	if iniSection.HasKey(key) {
		if doToLower {
			return strings.ToLower(iniSection.Key(key).String())
		} else {
			return iniSection.Key(key).String()
		}
	} else {
		return defaultValue
	}
}

func UnpackEtcIni(section ini.Section, allowBooleanKeys bool) (bool, *ini.File, string) {
	var iniFilePath = ""
	if get_from_section(section, "enabled", "false", true) != "true" {
		return false, nil, ""
	}
	iniFilePath = get_from_section(section, "ini_file_path", "", false)
	if iniFilePath == "" {
		return true, nil, ""
	}

	if allowBooleanKeys {
		var loadOptions ini.LoadOptions
		loadOptions.Insensitive = true
		loadOptions.AllowBooleanKeys = true
		cfg, err := ini.LoadSources(loadOptions, iniFilePath)

		if err != nil {
			log.Printf("Could not load (loadsources) ini file %s (%s)", iniFilePath, err)
			return true, nil, iniFilePath
		}
		return true, cfg, iniFilePath
	} else {
		cfg, err := ini.Load(iniFilePath)

		if err != nil {
			log.Printf("Could not load ini file %s (%s)", iniFilePath, err)
			return true, nil, iniFilePath
		}
		return true, cfg, iniFilePath
	}
}

func ReadEtcConfiguration(iniFileName string) (*ini.Section) {
	cfg, err := ini.Load(iniFileName)
	if err != nil {
		log.Fatalf("Ini Loader error: [%s] [%v]", iniFileName, err)
	}

	for _, sectionName := range cfg.SectionStrings() {
		if sectionName != "DEFAULT" {
			section, err := cfg.GetSection(sectionName)
			if err == nil {
				log.Printf("Loading %s config\n", sectionName)
				return section
			}
		}
	}

	return nil
}

func SaveIniFile(iniFile ini.File, iniFilePath string, testFileName string) {
	var test_output_folder = os.Getenv("TEST_OUTPUT_FOLDER")
	if test_output_folder == "" {
		log.Println("Writing config file to it's original location ", iniFilePath)
		iniFile.SaveTo(iniFilePath)
	} else {
		testFile := strings.Join([]string{test_output_folder, testFileName}, "/")
		log.Println("Writing config file to test file ", testFile)
		iniFile.SaveTo(testFile)
	}
}

func UpdateBoolKey(title string, section *ini.Section, key string, val bool) {
	log.Printf("%s config		%s = %t\n", title, key, val)
	convVal := strconv.FormatBool(val)
	convVal = strings.Title(convVal)
	section.NewKey(key, convVal)
}

func UpdateStringKey(title string, section *ini.Section, key string, val string) {
	log.Printf("%s config		%s = %s\n", title, key, val)
	section.NewKey(key, val)
}

func UpdateInt64Key(title string, section *ini.Section, key string, val int64) {
	log.Printf("%s config		%s = %d\n", title, key, val)
	section.NewKey(key, strconv.FormatInt(val, 10))
}

func EnsureDirExists(folderName string) {
		stat, err := os.Stat(folderName)
		if err != nil || !stat.IsDir() {
			mkdirerr := os.MkdirAll(folderName, 0777)
			if mkdirerr != nil {
				log.Printf("Could not create folder %s", folderName)
				return
			}
		}
}

func ReadLinesFromFile(full_file_path string) ([]string) {
	lines := []string{}
	stat, err := os.Stat(full_file_path)
	if err == nil && !stat.IsDir() {
		fh, err := os.Open(full_file_path)
		if err != nil {
			log.Printf("The following error occurred while trying to read %s : %s", full_file_path, err)
			return lines
		}

		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		fh.Close()
	} else {
		log.Printf("Unknown file read error: %s", err)
	}
	return lines
}

func WriteLinesToFile(full_file_path string, lines_to_write []string) {
	EnsureDirExists(path.Dir(full_file_path))

	fh, err := os.OpenFile(full_file_path, os.O_WRONLY, 0777)
	if err != nil {
		fh, err = os.Create(full_file_path)
		if err != nil {
			log.Printf("Could not write to %s: %s", full_file_path, err)
			return
		}
	}
	for _, line := range lines_to_write {
		fh.WriteString(fmt.Sprintf("%s\n", line))
	}
	fh.Close()
}
