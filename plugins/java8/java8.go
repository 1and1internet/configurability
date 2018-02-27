package main

import (
	"C"
	"log"
	"os"
	"path"
	"strings"

	"fmt"
	"math"

	"github.com/1and1internet/configurability/plugins"
	"github.com/go-ini/ini"
	yaml "gopkg.in/yaml.v2"
)

const OutputFileName = "/etc/configurability/custom/java_opts"

type Customisor interface {
	ApplyCustomisations()
}

// This is the entry point for the plugin...
func Customise(customisationFileContent []byte, section *ini.Section, configurationFileName string) bool {
	if OurConfigFileName(configurationFileName) {
		log.Println("Process as java8/yaml")
		allInfo := CustomisationInfo{
			JavaSwitches:    &Java8Data{},
			ConfJavaSection: *section,
		}
		if plugins.GetFromSection(*section, "enabled", "false", true) == "true" {
			if allInfo.LoadCustomConfig(customisationFileContent) {
				allInfo.GetMaxMemory()
				allInfo.ApplyCustomisations()
			}
		}
		return true
	}
	return false
}

func (allInfo *CustomisationInfo) LoadCustomConfig(customisationFileContent []byte) bool {
	err := yaml.Unmarshal(customisationFileContent, allInfo.JavaSwitches)
	if err != nil {
		log.Printf("error: %v", err)
		return false
	}
	return true
}

func getOutputFileName() string {
	var test_output_folder = os.Getenv("TEST_OUTPUT_FOLDER")
	if test_output_folder != "" {
		return path.Join(test_output_folder, OutputFileName)
	} else {
		return OutputFileName
	}
}

func (allInfo *CustomisationInfo) ApplyCustomisations() {
	var options []string

	allInfo.OptionNoclassgc(&options)
	allInfo.OptionAggressiveHeap(&options)
	allInfo.OptionDisableExplicitGC(&options)
	allInfo.OptionInitialHeapSize_Percent(&options)
	allInfo.OptionMaxGCPauseMillis(&options)
	allInfo.OptionMaxHeapSize_Percent(&options)
	allInfo.OptionMaxHeapFreeRatio(&options)
	allInfo.OptionMaxMetaspaceSize(&options)
	allInfo.OptionMaxNewSize_Percent(&options)
	allInfo.OptionMaxTenuringThreshold(&options)
	allInfo.OptionMinHeapFreeRatio(&options)
	allInfo.OptionNewRatio(&options)
	allInfo.OptionNewSize_Percent(&options)
	allInfo.OptionParallelGCThreads(&options)
	allInfo.OptionParallelRefProcEnabled(&options)

	line := strings.Join(options, " ")
	lines := []string{fmt.Sprintf("export CUSTOM_JAVA_OPTS=\"%s\"", line)}
	plugins.WriteLinesToFile(getOutputFileName(), lines)
}

func OurConfigFileName(configurationFileName string) bool {
	if configurationFileName == "configuration-java8.json" || configurationFileName == "configuration-java8.yaml" {
		return true
	}
	return false
}

func (allInfo *CustomisationInfo) GetMaxMemory() {
	allInfo.MaxMemoryBytes = 0
	cgroup_mem_limit_int, err := plugins.GetMaxMemoryOfContainer()
	if err == nil {
		if cgroup_mem_limit_int > 8589934592 {
			log.Printf("WARNING: Forcing 8G memory limit")
			cgroup_mem_limit_int = 8589934592
		}
		allInfo.MaxMemoryBytes = cgroup_mem_limit_int
	} else {
		log.Printf("ERROR getting container max memory %v", err)
		log.Print("WARNING: All memory options will be defaults")
	}
}

type CustomisationInfo struct {
	JavaSwitches         *Java8Data
	ConfJavaSection      ini.Section
	SourceConfigFilePath string
	MaxMemoryBytes       uint64
}

type Java8Data struct {
	// As defined in https://docs.oracle.com/javase/8/docs/technotes/tools/unix/java.html
	X  NonStandardOptions     `yaml:"X"`
	XX AdvancedRuntimeOptions `yaml:"XX"`
}

type NonStandardOptions struct {
	Noclassgc bool `yaml:"noclassgc"`
}

type AdvancedRuntimeOptions struct {
	AggressiveHeap          bool  `yaml:"AggressiveHeap"`
	DisableExplicitGC       bool  `yaml:"DisableExplicitGC"`
	InitialHeapSize_Percent int   `yaml:"InitialHeapSize_Percent"`
	MaxGCPauseMillis        int32 `yaml:"MaxGCPauseMillis"`
	MaxHeapSize_Percent     int   `yaml:"MaxHeapSize_Percent"`
	MaxHeapFreeRatio        int32 `yaml:"MaxHeapFreeRatio"`
	MaxMetaspaceSize        int32 `yaml:"MaxMetaspaceSize"`
	MaxNewSize_Percent      int   `yaml:"MaxNewSize_Percent"`
	MaxTenuringThreshold    int32 `yaml:"MaxTenuringThreshold"`
	MinHeapFreeRatio        int32 `yaml:"MinHeapFreeRatio"`
	NewRatio                int32 `yaml:"NewRatio"`
	NewSize_Percent         int   `yaml:"NewSize_Percent"`
	ParallelGCThreads       int32 `yaml:"ParallelGCThreads"`
	ParallelRefProcEnabled  bool  `yaml:"ParallelRefProcEnabled"`
}

func (allInfo *CustomisationInfo) OptionNoclassgc(options *[]string) {
	if allInfo.JavaSwitches.X.Noclassgc {
		*options = append(*options, "-Xnoclassgc")
	}
}

func (allInfo *CustomisationInfo) OptionAggressiveHeap(options *[]string) {
	if allInfo.JavaSwitches.XX.AggressiveHeap {
		*options = append(*options, "-XX:+AggressiveHeap")
	}
}

func (allInfo *CustomisationInfo) OptionDisableExplicitGC(options *[]string) {
	if allInfo.JavaSwitches.XX.DisableExplicitGC {
		*options = append(*options, "-XX:+DisableExplicitGC")
	}
}

func (allInfo *CustomisationInfo) GetJavaMemorySizeOptionString(
	percent int,
	defaultPercent int,
	thresholdPercent int,
	optstring string,
	options *[]string) {
	if percent <= thresholdPercent {
		percent = defaultPercent
	}
	if percent > 0 && percent <= 100 {
		size := uint64(float64(allInfo.MaxMemoryBytes) * (float64(percent) / 100.0))
		size_str := GetMemoryInMultiplesOf1024AsTidySuffixedString(size)
		option := fmt.Sprintf(optstring, size_str)
		*options = append(*options, option)
	}
}

func GetMemoryInMultiplesOf1024AsTidySuffixedString(size uint64) string {
	size = GetRoundedTo1024(size)
	size_str := fmt.Sprintf("%v", size)
	for _, suffix := range []string{"K", "M", "G"} {
		if size >= 1048576 {
			size = size / 1024
			size_str = fmt.Sprintf("%v%s", size, suffix)
		} else {
			break
		}
	}
	return size_str
}

func GetRoundedTo1024(size uint64) uint64 {
	remainder := math.Mod(float64(size), 1024.0)
	if remainder > 512.0 {
		return size + uint64(1024.0-remainder)
	} else if remainder > 0 {
		return size - uint64(remainder)
	}
	return size
}

func (allInfo *CustomisationInfo) OptionInitialHeapSize_Percent(options *[]string) {
	pc := allInfo.JavaSwitches.XX.InitialHeapSize_Percent
	optstring := "-XX:InitialHeapSize=%v"
	allInfo.GetJavaMemorySizeOptionString(pc, 0, -1, optstring, options)
}

func (allInfo *CustomisationInfo) OptionMaxGCPauseMillis(options *[]string) {
	value := allInfo.JavaSwitches.XX.MaxGCPauseMillis
	if value > 0 {
		option := fmt.Sprintf("-XX:MaxGCPauseMillis=%v", value)
		*options = append(*options, option)
	}
}

func (allInfo *CustomisationInfo) OptionMaxHeapSize_Percent(options *[]string) {
	pc := allInfo.JavaSwitches.XX.MaxHeapSize_Percent
	optstring := "-XX:MaxHeapSize=%v"
	allInfo.GetJavaMemorySizeOptionString(pc, 100, 0, optstring, options)
}

func (allInfo *CustomisationInfo) OptionMaxHeapFreeRatio(options *[]string) {
	value := allInfo.JavaSwitches.XX.MaxHeapFreeRatio
	if value > 0 {
		option := fmt.Sprintf("-XX:MaxHeapFreeRatio=%v", value)
		*options = append(*options, option)
	}
}

func (allInfo *CustomisationInfo) OptionMaxMetaspaceSize(options *[]string) {
	value := allInfo.JavaSwitches.XX.MaxMetaspaceSize
	if value > 0 {
		option := fmt.Sprintf("-XX:MaxMetaspaceSize=%v", value)
		*options = append(*options, option)
	}
}

func (allInfo *CustomisationInfo) OptionMaxNewSize_Percent(options *[]string) {
	pc := allInfo.JavaSwitches.XX.MaxNewSize_Percent
	optstring := "-XX:MaxNewSize=%v"
	allInfo.GetJavaMemorySizeOptionString(pc, 50, 0, optstring, options)
}

func (allInfo *CustomisationInfo) OptionMaxTenuringThreshold(options *[]string) {
	value := allInfo.JavaSwitches.XX.MaxTenuringThreshold
	if value > 0 {
		option := fmt.Sprintf("-XX:MaxTenuringThreshold=%v", value)
		*options = append(*options, option)
	}
}

func (allInfo *CustomisationInfo) OptionMinHeapFreeRatio(options *[]string) {
	value := allInfo.JavaSwitches.XX.MinHeapFreeRatio
	if value > 0 {
		option := fmt.Sprintf("-XX:MinHeapFreeRatio=%v", value)
		*options = append(*options, option)
	}
}

func (allInfo *CustomisationInfo) OptionNewRatio(options *[]string) {
	value := allInfo.JavaSwitches.XX.NewRatio
	if value > 0 {
		option := fmt.Sprintf("-XX:NewRatio=%v", value)
		*options = append(*options, option)
	}
}

func (allInfo *CustomisationInfo) OptionNewSize_Percent(options *[]string) {
	pc := allInfo.JavaSwitches.XX.NewSize_Percent
	optstring := "-XX:NewSize=%v"
	allInfo.GetJavaMemorySizeOptionString(pc, 50, 0, optstring, options)
}

func (allInfo *CustomisationInfo) OptionParallelGCThreads(options *[]string) {
	value := allInfo.JavaSwitches.XX.ParallelGCThreads
	if value > 0 {
		option := fmt.Sprintf("-XX:ParallelGCThreads=%v", value)
		*options = append(*options, option)
	}
}

func (allInfo *CustomisationInfo) OptionParallelRefProcEnabled(options *[]string) {
	if allInfo.JavaSwitches.XX.ParallelRefProcEnabled {
		*options = append(*options, "-XX:+ParallelRefProcEnabled")
	}
}
