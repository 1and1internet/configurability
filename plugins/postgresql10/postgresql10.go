package main

import (
	"C"
	"log"
	"os"
	"path"
	"strings"

	"github.com/1and1internet/configurability/plugins"
	"github.com/go-ini/ini"
	yaml "gopkg.in/yaml.v2"
)
import (
	"fmt"
	"strconv"
)

const OutputFileName = "/var/lib/postgresql/10/main/postgresql.conf"

type Customisor interface {
	ApplyCustomisations()
}

// This is the entry point for the plugin...
func Customise(customisationFileContent []byte, section *ini.Section, configurationFileName string) bool {
	if OurConfigFileName(configurationFileName) {
		log.Println("Process as PostgreSql10/yaml")
		allInfo := CustomisationInfo{
			PostgreSqlRequestedConfig:  &RequestedConfig{},
			ConfPostgreSqlSection:      section,
			PostgresSqlDotConfFilename: plugins.GetFromSection(*section, "ini_file_path", "", false),
		}
		if plugins.GetFromSection(*section, "enabled", "false", true) == "true" {
			if allInfo.LoadCustomConfig(customisationFileContent) && allInfo.LoadCurrentConfig() {
				allInfo.SetMaxMemory()
				allInfo.ApplyCustomisations()
			}
		}
		return true
	}
	return false
}

func (allInfo *CustomisationInfo) LoadCustomConfig(customisationFileContent []byte) bool {
	err := yaml.Unmarshal(customisationFileContent, allInfo.PostgreSqlRequestedConfig)
	if err != nil {
		log.Printf("error: %v", err)
		return false
	}
	return true
}

func (allInfo *CustomisationInfo) LoadCurrentConfig() bool {
	allInfo.PostgresSqlDotConfLines = plugins.ReadLinesFromFile(allInfo.PostgresSqlDotConfFilename)
	if len(allInfo.PostgresSqlDotConfLines) == 0 {
		return false
	}
	allInfo.ParsePostgresSqlDotConf()
	return true
}

func (allInfo *CustomisationInfo) ParsePostgresSqlDotConf() {
	allInfo.ParsedConfLineMap = make(map[string]*ConfigLine)
	for _, line := range allInfo.PostgresSqlDotConfLines {
		aConfline := ConfigLine{
			OrigLine:         line,
			UseOrig:          true,
			Key:              "",
			Value:            "",
			PostValueComment: "",
		}
		subline := ""
		if len(line) > 0 {
			doParse := true
			if line[0] == '#' {
				subline = line[1:]
			} else {
				//aConfline.UseOrig = false
				subline = line
				if line[0] == ' ' || line[0] == '\t' {
					doParse = false
				}
			}
			if doParse {
				keyvalue := strings.Split(subline, "=")
				if len(keyvalue) > 1 {
					aConfline.Key = strings.Trim(keyvalue[0], " ")
					valuecomment := strings.Split(keyvalue[1], "#")
					if len(valuecomment) > 1 {
						aConfline.Value = strings.Trim(valuecomment[0], " '\t")
						aConfline.PostValueComment = valuecomment[1]
					} else {
						aConfline.Value = strings.Trim(keyvalue[1], " '\t")
						aConfline.PostValueComment = ""
					}
				}
			}
		}
		allInfo.ParsedConfLines = append(allInfo.ParsedConfLines, &aConfline)
		if aConfline.Key != "" {
			allInfo.ParsedConfLineMap[aConfline.Key] = &aConfline
		}
	}
}

func getOutputFileName() string {
	var test_output_folder = os.Getenv("TEST_OUTPUT_FOLDER")
	if test_output_folder != "" {
		return path.Join(test_output_folder, OutputFileName)
	} else {
		return OutputFileName
	}
}

func OurConfigFileName(configurationFileName string) bool {
	if configurationFileName == "configuration-postgresql10.json" || configurationFileName == "configuration-postgresql10.yaml" {
		return true
	}
	return false
}

func (allInfo *CustomisationInfo) SetMaxMemory() {
	allInfo.MaxMemory = plugins.GetMemoryValue(
		plugins.GetMaxMemoryOfContainerAsString("17179869184"),
	)
	allInfo.MaxMemory.MemStrToMemValue()
}

func (allInfo *CustomisationInfo) ApplyCustomisations() {
	allInfo.MaxConnections()
	allInfo.SuperuserReservedConnections()
	allInfo.SharedBuffers()
	allInfo.HugePages()
	allInfo.MaxPreparedTransactions()
	allInfo.WorkMem()
	allInfo.MaintenanceWorkMem()
	allInfo.ReplacementSortTuples()
	allInfo.MaxStackDepth()
	allInfo.VacuumCostDelay()
	allInfo.VacuumCostPageHit()
	allInfo.VacuumCostPageMiss()
	allInfo.VacuumCostPageDirty()
	allInfo.VacuumCostLimit()
	allInfo.BgwriterDelay()
	allInfo.BgwriterLruMaxpages()
	allInfo.BgwriterLruMultiplier()
	allInfo.WalLevel()
	allInfo.SynchronousCommit()
	allInfo.WalSyncMethod()
	allInfo.WalLogHints()
	allInfo.WalCompression()
	allInfo.WalWriterDelay()
	allInfo.WalWriterFlushAfter()
	allInfo.CommitDelay()
	allInfo.CommitSiblings()
	allInfo.CheckpointTimeout()
	allInfo.CheckpointCompletionTarget()
	allInfo.CheckpointFlushAfter()
	allInfo.CheckpointWarning()
	allInfo.MaxWalSize()
	allInfo.MinWalSize()
	allInfo.ArchiveMode()
	allInfo.ArchiveCommand()
	allInfo.ArchiveTimeout()
	allInfo.EnableBitmapscan()
	allInfo.EnableHashagg()
	allInfo.EnableHashjoin()
	allInfo.EnableIndexscan()
	allInfo.EnableIndexonlyscan()
	allInfo.EnableMaterial()
	allInfo.EnableMergejoin()
	allInfo.EnableNestloop()
	allInfo.EnableSeqscan()
	allInfo.EnableSort()
	allInfo.EnableTidscan()
	allInfo.LogDestination()
	allInfo.LoggingCollector()
	allInfo.ClientMinMessages()
	allInfo.LogMinMessages()
	allInfo.LogMinErrorStatement()
	allInfo.LogMinDurationStatement()
	allInfo.DebugPrintParse()
	allInfo.DebugPrintRewritten()
	allInfo.DebugPrintPlan()
	allInfo.DebugPrettyPrint()
	allInfo.LogCheckpoints()
	allInfo.LogConnections()
	allInfo.LogDisconnections()
	allInfo.LogDuration()
	allInfo.LogErrorVerbosity()
	allInfo.LogHostname()
	//allInfo.LogLinePrefix()
	allInfo.LogLockWaits()
	allInfo.LogStatement()
	allInfo.LogTempFiles()
	allInfo.TrackActivities()
	allInfo.TrackCounts()
	allInfo.TrackIoTiming()
	allInfo.TrackFunctions()
	allInfo.TrackActivityQuerySize()
	allInfo.LogParserStats()
	allInfo.LogPlannerStats()
	allInfo.LogExecutorStats()
	allInfo.LogStatementStats()
	allInfo.Autovacuum()
	allInfo.LogAutovacuumMinDuration()
	allInfo.AutovacuumMaxWorkers()
	allInfo.AutovacuumNaptime()
	allInfo.AutovacuumVacuumThreshold()
	allInfo.AutovacuumAnalyzeThreshold()
	allInfo.AutovacuumVacuumScaleFactor()
	allInfo.AutovacuumAnalyzeScaleFactor()
	allInfo.AutovacuumFreezeMaxAge()
	allInfo.AutovacuumMultixactFreezeMaxAge()
	allInfo.AutovacuumVacuumCostDelay()
	allInfo.AutovacuumVacuumCostLimit()

	lines := []string{}
	for _, confline := range allInfo.ParsedConfLines {
		if confline.UseOrig {
			lines = append(lines, confline.OrigLine)
		} else {
			line := fmt.Sprintf("%s = %s\t#%s", confline.Key, confline.Value, confline.PostValueComment)
			lines = append(lines, line)
		}
	}
	outfile := getOutputFileName()
	log.Printf("Writing %d lines to %s", len(lines), outfile)
	plugins.WriteLinesToFile(outfile, lines)
}

func (confline *ConfigLine) SetIntVal(value, defaultValue, min, max int) {
	if value < min || value > max {
		value = defaultValue
	}
	valuestr := strconv.Itoa(value)
	if valuestr != confline.Value {
		confline.Value = valuestr
		confline.UseOrig = false
	}
}

func (confline *ConfigLine) SetFloatVal(value, defaultValue, min, max string) {
	var floatValue, floatMin, floatMax float64
	var err error

	floatValue, err = strconv.ParseFloat(value, 32)
	if err != nil {
		log.Printf("Error converting to float: %s", value)
	}

	floatMin, err = strconv.ParseFloat(min, 32)
	if err != nil {
		log.Printf("Error converting to float: %s", min)
	}

	floatMax, err = strconv.ParseFloat(max, 32)
	if err != nil {
		log.Printf("Error converting to float: %s", max)
	}

	if floatValue < floatMin || floatValue > floatMax {
		value = defaultValue
	}

	if confline.Value != value {
		confline.Value = value
		confline.UseOrig = false
	}
}

func (confline *ConfigLine) SetStrVal(value, defaultValue string) {
	if value == "" {
		value = defaultValue
	}
	if value != confline.Value {
		confline.Value = value
		confline.UseOrig = false
	}
}

func (confline *ConfigLine) SetTimeVal(value, defaultValue, min, max string) {
	timeValue := plugins.GetTimeValue(value)
	timeMin := plugins.GetTimeValue(min)
	timeMax := plugins.GetTimeValue(max)

	if timeValue.Error != nil || timeMin.Error != nil || timeMax.Error != nil {
		return
	}
	if timeValue.LessThan(timeMin) || timeMax.LessThan(timeValue) {
		value = defaultValue
	}
	if confline.Value != value {
		confline.Value = value
		confline.UseOrig = false
	}
}

func (confline *ConfigLine) SetMemVal(value, defaultValue, min, max string, systemMax plugins.MemValue) {
	memValue := plugins.GetMemoryValue(value)
	minValue := plugins.GetMemoryValue(min)
	maxValue := plugins.GetMemoryValue(max)

	if memValue.Error != nil || minValue.Error != nil || maxValue.Error != nil {
		return
	}

	if memValue.LessThan(minValue) || maxValue.LessThan(memValue) {
		value = defaultValue
	}

	if value != confline.Value {
		confline.Value = memValue.CorrectStrValue
		confline.UseOrig = false
	}
}

func (allInfo *CustomisationInfo) MaxConnections() {
	allInfo.ParsedConfLineMap["max_connections"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.MaxConnections, 100, 1, 100000,
	)
}

func (allInfo *CustomisationInfo) SuperuserReservedConnections() {
	allInfo.ParsedConfLineMap["superuser_reserved_connections"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.SuperuserReservedConnections,
		3,
		1,
		99999,
	)
}

func (allInfo *CustomisationInfo) SharedBuffers() {
	allInfo.ParsedConfLineMap["shared_buffers"].SetMemVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.SharedBuffers,
		"128MB",
		"128kB",
		"1GB",
		*allInfo.MaxMemory,
	)
}

func (allInfo *CustomisationInfo) HugePages() {
	allInfo.ParsedConfLineMap["huge_pages"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.HugePages, "try",
	)
}

func (allInfo *CustomisationInfo) MaxPreparedTransactions() {
	allInfo.ParsedConfLineMap["max_prepared_transactions"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.MaxPreparedTransactions,
		0,
		0,
		9999999,
	)
}

func (allInfo *CustomisationInfo) WorkMem() {
	allInfo.ParsedConfLineMap["work_mem"].SetMemVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.WorkMem,
		"4MB",
		"512kB",
		"50MB",
		*allInfo.MaxMemory,
	)
}

func (allInfo *CustomisationInfo) MaintenanceWorkMem() {
	allInfo.ParsedConfLineMap["maintenance_work_mem"].SetMemVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.MaintenanceWorkMem,
		"64MB",
		"512kB",
		"1GB",
		*allInfo.MaxMemory,
	)
}

func (allInfo *CustomisationInfo) ReplacementSortTuples() {
	allInfo.ParsedConfLineMap["replacement_sort_tuples"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.ReplacementSortTuples,
		150000,
		0,
		9999999,
	)
}

func (allInfo *CustomisationInfo) MaxStackDepth() {
	allInfo.ParsedConfLineMap["max_stack_depth"].SetMemVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.MaxStackDepth,
		"2MB",
		"512kB",
		"1GB",
		*allInfo.MaxMemory,
	)
}

func (allInfo *CustomisationInfo) VacuumCostDelay() {
	allInfo.ParsedConfLineMap["vacuum_cost_delay"].SetTimeVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.VacuumCostDelay,
		"0",
		"0",
		"100ms",
	)
}

func (allInfo *CustomisationInfo) VacuumCostPageHit() {
	allInfo.ParsedConfLineMap["vacuum_cost_page_hit"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.VacuumCostPageHit,
		1,
		0,
		10000,
	)
}

func (allInfo *CustomisationInfo) VacuumCostPageMiss() {
	allInfo.ParsedConfLineMap["vacuum_cost_page_miss"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.VacuumCostPageMiss,
		10,
		0,
		10000,
	)
}

func (allInfo *CustomisationInfo) VacuumCostPageDirty() {
	allInfo.ParsedConfLineMap["vacuum_cost_page_dirty"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.VacuumCostPageDirty,
		20,
		0,
		10000,
	)
}

func (allInfo *CustomisationInfo) VacuumCostLimit() {
	allInfo.ParsedConfLineMap["vacuum_cost_limit"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.VacuumCostLimit,
		200,
		0,
		10000,
	)
}

func (allInfo *CustomisationInfo) BgwriterDelay() {
	allInfo.ParsedConfLineMap["bgwriter_delay"].SetTimeVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.BgwriterDelay,
		"200ms",
		"10ms",
		"10000ms",
	)
}

func (allInfo *CustomisationInfo) BgwriterLruMaxpages() {
	allInfo.ParsedConfLineMap["bgwriter_lru_maxpages"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.BgwriterLruMaxpages,
		100,
		0,
		1000,
	)
}

func (allInfo *CustomisationInfo) BgwriterLruMultiplier() {
	allInfo.ParsedConfLineMap["bgwriter_lru_multiplier"].SetFloatVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.BgwriterLruMultiplier,
		"2.0",
		"0.0",
		"10.0",
	)
}

func (allInfo *CustomisationInfo) WalLevel() {
	allInfo.ParsedConfLineMap["wal_level"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.WalLevel, "replica",
	)
}

func (allInfo *CustomisationInfo) SynchronousCommit() {
	allInfo.ParsedConfLineMap["synchronous_commit"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.SynchronousCommit, "on",
	)
}

func (allInfo *CustomisationInfo) WalSyncMethod() {
	allInfo.ParsedConfLineMap["wal_sync_method"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.WalSyncMethod, "fdatasync",
	)
}

func (allInfo *CustomisationInfo) WalLogHints() {
	allInfo.ParsedConfLineMap["wal_log_hints"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.WalLogHints, "off",
	)
}

func (allInfo *CustomisationInfo) WalCompression() {
	allInfo.ParsedConfLineMap["wal_compression"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.WalCompression, "off",
	)
}

func (allInfo *CustomisationInfo) WalWriterDelay() {
	allInfo.ParsedConfLineMap["wal_writer_delay"].SetTimeVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.WalWriterDelay,
		"200ms",
		"1ms",
		"10000ms",
	)
}

func (allInfo *CustomisationInfo) WalWriterFlushAfter() {
	allInfo.ParsedConfLineMap["wal_writer_flush_after"].SetMemVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.WalWriterFlushAfter,
		"1MB",
		"0",
		"100MB",
		*allInfo.MaxMemory,
	)
}

func (allInfo *CustomisationInfo) CommitDelay() {
	allInfo.ParsedConfLineMap["commit_delay"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.CommitDelay,
		0,
		0,
		100000,
	)
}

func (allInfo *CustomisationInfo) CommitSiblings() {
	allInfo.ParsedConfLineMap["commit_siblings"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.CommitSiblings,
		5,
		1,
		1000,
	)
}

func (allInfo *CustomisationInfo) CheckpointTimeout() {
	allInfo.ParsedConfLineMap["checkpoint_timeout"].SetTimeVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.CheckpointTimeout,
		"5min",
		"30s",
		"1d",
	)
}

func (allInfo *CustomisationInfo) CheckpointCompletionTarget() {
	allInfo.ParsedConfLineMap["checkpoint_completion_target"].SetFloatVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.CheckpointCompletionTarget,
		"0.5",
		"0.0",
		"1.0",
	)
}

func (allInfo *CustomisationInfo) CheckpointFlushAfter() {
	allInfo.ParsedConfLineMap["checkpoint_flush_after"].SetMemVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.CheckpointFlushAfter,
		"256kB",
		"0",
		"500MB",
		*allInfo.MaxMemory,
	)
}

func (allInfo *CustomisationInfo) CheckpointWarning() {
	allInfo.ParsedConfLineMap["checkpoint_warning"].SetTimeVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.CheckpointWarning,
		"30s",
		"0",
		"10d",
	)
}

func (allInfo *CustomisationInfo) MaxWalSize() {
	allInfo.ParsedConfLineMap["max_wal_size"].SetMemVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.MaxWalSize,
		"1GB",
		"0",
		"500GB",
		*allInfo.MaxMemory,
	)
}

func (allInfo *CustomisationInfo) MinWalSize() {
	allInfo.ParsedConfLineMap["min_wal_size"].SetMemVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.MinWalSize,
		"80MB",
		"0",
		"500GB",
		*allInfo.MaxMemory,
	)
}

func (allInfo *CustomisationInfo) ArchiveMode() {
	allInfo.ParsedConfLineMap["archive_mode"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.ArchiveMode, "off",
	)
}

func (allInfo *CustomisationInfo) ArchiveCommand() {
	allInfo.ParsedConfLineMap["archive_command"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.ArchiveCommand, "",
	)
}

func (allInfo *CustomisationInfo) ArchiveTimeout() {
	allInfo.ParsedConfLineMap["archive_timeout"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.ArchiveTimeout,
		0,
		0,
		100000,
	)
}

func (allInfo *CustomisationInfo) EnableBitmapscan() {
	allInfo.ParsedConfLineMap["enable_bitmapscan"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableBitmapscan, "on",
	)
}

func (allInfo *CustomisationInfo) EnableHashagg() {
	allInfo.ParsedConfLineMap["enable_hashagg"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableHashagg, "on",
	)
}

func (allInfo *CustomisationInfo) EnableHashjoin() {
	allInfo.ParsedConfLineMap["enable_hashjoin"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableHashjoin, "on",
	)
}

func (allInfo *CustomisationInfo) EnableIndexscan() {
	allInfo.ParsedConfLineMap["enable_indexscan"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableIndexscan, "on",
	)
}

func (allInfo *CustomisationInfo) EnableIndexonlyscan() {
	allInfo.ParsedConfLineMap["enable_indexonlyscan"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableIndexonlyscan, "on",
	)
}

func (allInfo *CustomisationInfo) EnableMaterial() {
	allInfo.ParsedConfLineMap["enable_material"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableMaterial, "on",
	)
}

func (allInfo *CustomisationInfo) EnableMergejoin() {
	allInfo.ParsedConfLineMap["enable_mergejoin"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableMergejoin, "on",
	)
}

func (allInfo *CustomisationInfo) EnableNestloop() {
	allInfo.ParsedConfLineMap["enable_nestloop"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableNestloop, "on",
	)
}

func (allInfo *CustomisationInfo) EnableSeqscan() {
	allInfo.ParsedConfLineMap["enable_seqscan"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableSeqscan, "on",
	)
}

func (allInfo *CustomisationInfo) EnableSort() {
	allInfo.ParsedConfLineMap["enable_sort"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableSort, "on",
	)
}

func (allInfo *CustomisationInfo) EnableTidscan() {
	allInfo.ParsedConfLineMap["enable_tidscan"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.EnableTidscan, "on",
	)
}

func (allInfo *CustomisationInfo) LogDestination() {
	allInfo.ParsedConfLineMap["log_destination"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogDestination, "stderr",
	)
}

func (allInfo *CustomisationInfo) LoggingCollector() {
	allInfo.ParsedConfLineMap["logging_collector"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LoggingCollector, "off",
	)
}

func (allInfo *CustomisationInfo) ClientMinMessages() {
	allInfo.ParsedConfLineMap["client_min_messages"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.ClientMinMessages, "NOTICE",
	)
}

func (allInfo *CustomisationInfo) LogMinMessages() {
	allInfo.ParsedConfLineMap["log_min_messages"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogMinMessages, "WARNING",
	)
}

func (allInfo *CustomisationInfo) LogMinErrorStatement() {
	allInfo.ParsedConfLineMap["log_min_error_statement"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogMinErrorStatement, "ERROR",
	)
}

func (allInfo *CustomisationInfo) LogMinDurationStatement() {
	allInfo.ParsedConfLineMap["log_min_duration_statement"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogMinDurationStatement,
		-1,
		-1,
		10000,
	)
}

func (allInfo *CustomisationInfo) DebugPrintParse() {
	allInfo.ParsedConfLineMap["debug_print_parse"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.DebugPrintParse, "off",
	)
}

func (allInfo *CustomisationInfo) DebugPrintRewritten() {
	allInfo.ParsedConfLineMap["debug_print_rewritten"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.DebugPrintRewritten, "off",
	)
}

func (allInfo *CustomisationInfo) DebugPrintPlan() {
	allInfo.ParsedConfLineMap["debug_print_plan"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.DebugPrintPlan, "off",
	)
}

func (allInfo *CustomisationInfo) DebugPrettyPrint() {
	allInfo.ParsedConfLineMap["debug_pretty_print"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.DebugPrettyPrint, "on",
	)
}

func (allInfo *CustomisationInfo) LogCheckpoints() {
	allInfo.ParsedConfLineMap["log_checkpoints"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogCheckpoints, "off",
	)
}

func (allInfo *CustomisationInfo) LogConnections() {
	allInfo.ParsedConfLineMap["log_connections"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogConnections, "off",
	)
}

func (allInfo *CustomisationInfo) LogDisconnections() {
	allInfo.ParsedConfLineMap["log_disconnections"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogDisconnections, "off",
	)
}

func (allInfo *CustomisationInfo) LogDuration() {
	allInfo.ParsedConfLineMap["log_duration"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogDuration, "off",
	)
}

func (allInfo *CustomisationInfo) LogErrorVerbosity() {
	allInfo.ParsedConfLineMap["log_error_verbosity"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogErrorVerbosity, "default",
	)
}

func (allInfo *CustomisationInfo) LogHostname() {
	allInfo.ParsedConfLineMap["log_hostname"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogHostname, "off",
	)
}

//func (allInfo *CustomisationInfo) LogLinePrefix() {
//	allInfo.ParsedConfLineMap["log_line_prefix"].SetStrVal(
//		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogLinePrefix, "%m [%p] ",
//	)
//}

func (allInfo *CustomisationInfo) LogLockWaits() {
	allInfo.ParsedConfLineMap["log_lock_waits"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogLockWaits, "off",
	)
}

func (allInfo *CustomisationInfo) LogStatement() {
	allInfo.ParsedConfLineMap["log_statement"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogStatement, "none",
	)
}

func (allInfo *CustomisationInfo) LogTempFiles() {
	allInfo.ParsedConfLineMap["log_temp_files"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogTempFiles,
		-1,
		-1,
		10000,
	)
}

func (allInfo *CustomisationInfo) TrackActivities() {
	allInfo.ParsedConfLineMap["track_activities"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.TrackActivities, "on",
	)
}

func (allInfo *CustomisationInfo) TrackCounts() {
	allInfo.ParsedConfLineMap["track_counts"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.TrackCounts, "on",
	)
}

func (allInfo *CustomisationInfo) TrackIoTiming() {
	allInfo.ParsedConfLineMap["track_io_timing"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.TrackIoTiming, "off",
	)
}

func (allInfo *CustomisationInfo) TrackFunctions() {
	allInfo.ParsedConfLineMap["track_functions"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.TrackFunctions, "none",
	)
}

func (allInfo *CustomisationInfo) TrackActivityQuerySize() {
	allInfo.ParsedConfLineMap["track_activity_query_size"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.TrackActivityQuerySize,
		1024,
		0,
		10240,
	)
}

func (allInfo *CustomisationInfo) LogParserStats() {
	allInfo.ParsedConfLineMap["log_parser_stats"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogParserStats, "off",
	)
}

func (allInfo *CustomisationInfo) LogPlannerStats() {
	allInfo.ParsedConfLineMap["log_planner_stats"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogPlannerStats, "off",
	)
}

func (allInfo *CustomisationInfo) LogExecutorStats() {
	allInfo.ParsedConfLineMap["log_executor_stats"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogExecutorStats, "off",
	)
}

func (allInfo *CustomisationInfo) LogStatementStats() {
	allInfo.ParsedConfLineMap["log_statement_stats"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogStatementStats, "off",
	)
}

func (allInfo *CustomisationInfo) Autovacuum() {
	allInfo.ParsedConfLineMap["autovacuum"].SetStrVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.Autovacuum, "on",
	)
}

func (allInfo *CustomisationInfo) LogAutovacuumMinDuration() {
	allInfo.ParsedConfLineMap["log_autovacuum_min_duration"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.LogAutovacuumMinDuration,
		-1,
		-1,
		99999,
	)
}

func (allInfo *CustomisationInfo) AutovacuumMaxWorkers() {
	allInfo.ParsedConfLineMap["autovacuum_max_workers"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumMaxWorkers,
		3,
		1,
		99999,
	)
}

func (allInfo *CustomisationInfo) AutovacuumNaptime() {
	allInfo.ParsedConfLineMap["autovacuum_naptime"].SetTimeVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumNaptime,
		"1min",
		"1s",
		"1d",
	)
}

func (allInfo *CustomisationInfo) AutovacuumVacuumThreshold() {
	allInfo.ParsedConfLineMap["autovacuum_vacuum_threshold"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumVacuumThreshold,
		50,
		1,
		99999,
	)
}

func (allInfo *CustomisationInfo) AutovacuumAnalyzeThreshold() {
	allInfo.ParsedConfLineMap["autovacuum_analyze_threshold"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumAnalyzeThreshold,
		50,
		1,
		99999,
	)
}

func (allInfo *CustomisationInfo) AutovacuumVacuumScaleFactor() {
	allInfo.ParsedConfLineMap["autovacuum_vacuum_scale_factor"].SetFloatVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumVacuumScaleFactor,
		"0.2",
		"0.0",
		"1.0",
	)
}

func (allInfo *CustomisationInfo) AutovacuumAnalyzeScaleFactor() {
	allInfo.ParsedConfLineMap["autovacuum_analyze_scale_factor"].SetFloatVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumAnalyzeScaleFactor,
		"0.1",
		"0.0",
		"1.0",
	)
}

func (allInfo *CustomisationInfo) AutovacuumFreezeMaxAge() {
	allInfo.ParsedConfLineMap["autovacuum_freeze_max_age"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumFreezeMaxAge,
		200000000,
		1,
		9000000000,
	)
}

func (allInfo *CustomisationInfo) AutovacuumMultixactFreezeMaxAge() {
	allInfo.ParsedConfLineMap["autovacuum_multixact_freeze_max_age"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumMultixactFreezeMaxAge,
		400000000,
		1,
		9000000000,
	)
}

func (allInfo *CustomisationInfo) AutovacuumVacuumCostDelay() {
	allInfo.ParsedConfLineMap["autovacuum_vacuum_cost_delay"].SetTimeVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumVacuumCostDelay,
		"20ms",
		"-1",
		"1d",
	)
}

func (allInfo *CustomisationInfo) AutovacuumVacuumCostLimit() {
	allInfo.ParsedConfLineMap["autovacuum_vacuum_cost_limit"].SetIntVal(
		allInfo.PostgreSqlRequestedConfig.PostgreSQLConf.AutovacuumVacuumCostLimit,
		-1,
		-1,
		9999999,
	)
}

type CustomisationInfo struct {
	PostgreSqlRequestedConfig  *RequestedConfig
	ConfPostgreSqlSection      *ini.Section
	SourceConfigFilePath       string
	MaxMemory                  *plugins.MemValue
	PostgresSqlDotConfFilename string
	PostgresSqlDotConfLines    []string
	ParsedConfLines            []*ConfigLine
	ParsedConfLineMap          map[string]*ConfigLine
}

type RequestedConfig struct {
	PostgreSQLConf PostgreSqlConfSettings `yaml:"postgresqlconf"`
}

type PostgreSqlConfSettings struct {
	MaxConnections               int    `yaml:"max_connections"`
	SuperuserReservedConnections int    `yaml:"superuser_reserved_connections"`
	SharedBuffers                string `yaml:"shared_buffers"`
	HugePages                    string `yaml:"huge_pages"`
	MaxPreparedTransactions      int    `yaml:"max_prepared_transactions"`
	WorkMem                      string `yaml:"work_mem"`
	MaintenanceWorkMem           string `yaml:"maintenance_work_mem"`
	ReplacementSortTuples        int    `yaml:"replacement_sort_tuples"`
	MaxStackDepth                string `yaml:"max_stack_depth"`
	VacuumCostDelay              string `yaml:"vacuum_cost_delay"`
	VacuumCostPageHit            int    `yaml:"vacuum_cost_page_hit"`
	VacuumCostPageMiss           int    `yaml:"vacuum_cost_page_miss"`
	VacuumCostPageDirty          int    `yaml:"vacuum_cost_page_dirty"`
	VacuumCostLimit              int    `yaml:"vacuum_cost_limit"`
	BgwriterDelay                string `yaml:"bgwriter_delay"`
	BgwriterLruMaxpages          int    `yaml:"bgwriter_lru_maxpages"`
	BgwriterLruMultiplier        string `yaml:"bgwriter_lru_multiplier"`
	WalLevel                     string `yaml:"wal_level"`
	SynchronousCommit            string `yaml:"synchronous_commit"`
	WalSyncMethod                string `yaml:"wal_sync_method"`
	WalLogHints                  string `yaml:"wal_log_hints"`
	WalCompression               string `yaml:"wal_compression"`
	WalWriterDelay               string `yaml:"wal_writer_delay"`
	WalWriterFlushAfter          string `yaml:"wal_writer_flush_after"`
	CommitDelay                  int    `yaml:"commit_delay"`
	CommitSiblings               int    `yaml:"commit_siblings"`
	CheckpointTimeout            string `yaml:"checkpoint_timeout"`
	CheckpointCompletionTarget   string `yaml:"checkpoint_completion_target"`
	CheckpointFlushAfter         string `yaml:"checkpoint_flush_after"`
	CheckpointWarning            string `yaml:"checkpoint_warning"`
	MaxWalSize                   string `yaml:"max_wal_size"`
	MinWalSize                   string `yaml:"min_wal_size"`
	ArchiveMode                  string `yaml:"archive_mode"`
	ArchiveCommand               string `yaml:"archive_command"`
	ArchiveTimeout               int    `yaml:"archive_timeout"`
	EnableBitmapscan             string `yaml:"enable_bitmapscan"`
	EnableHashagg                string `yaml:"enable_hashagg"`
	EnableHashjoin               string `yaml:"enable_hashjoin"`
	EnableIndexscan              string `yaml:"enable_indexscan"`
	EnableIndexonlyscan          string `yaml:"enable_indexonlyscan"`
	EnableMaterial               string `yaml:"enable_material"`
	EnableMergejoin              string `yaml:"enable_mergejoin"`
	EnableNestloop               string `yaml:"enable_nestloop"`
	EnableSeqscan                string `yaml:"enable_seqscan"`
	EnableSort                   string `yaml:"enable_sort"`
	EnableTidscan                string `yaml:"enable_tidscan"`
	LogDestination               string `yaml:"log_destination"`
	LoggingCollector             string `yaml:"logging_collector"`
	ClientMinMessages            string `yaml:"client_min_messages"`
	LogMinMessages               string `yaml:"log_min_messages"`
	LogMinErrorStatement         string `yaml:"log_min_error_statement"`
	LogMinDurationStatement      int    `yaml:"log_min_duration_statement"`
	DebugPrintParse              string `yaml:"debug_print_parse"`
	DebugPrintRewritten          string `yaml:"debug_print_rewritten"`
	DebugPrintPlan               string `yaml:"debug_print_plan"`
	DebugPrettyPrint             string `yaml:"debug_pretty_print"`
	LogCheckpoints               string `yaml:"log_checkpoints"`
	LogConnections               string `yaml:"log_connections"`
	LogDisconnections            string `yaml:"log_disconnections"`
	LogDuration                  string `yaml:"log_duration"`
	LogErrorVerbosity            string `yaml:"log_error_verbosity"`
	LogHostname                  string `yaml:"log_hostname"`
	//LogLinePrefix                   string `yaml:"log_line_prefix"`
	LogLockWaits                    string `yaml:"log_lock_waits"`
	LogStatement                    string `yaml:"log_statement"`
	LogTempFiles                    int    `yaml:"log_temp_files"`
	TrackActivities                 string `yaml:"track_activities"`
	TrackCounts                     string `yaml:"track_counts"`
	TrackIoTiming                   string `yaml:"track_io_timing"`
	TrackFunctions                  string `yaml:"track_functions"`
	TrackActivityQuerySize          int    `yaml:"track_activity_query_size"`
	LogParserStats                  string `yaml:"log_parser_stats"`
	LogPlannerStats                 string `yaml:"log_planner_stats"`
	LogExecutorStats                string `yaml:"log_executor_stats"`
	LogStatementStats               string `yaml:"log_statement_stats"`
	Autovacuum                      string `yaml:"autovacuum"`
	LogAutovacuumMinDuration        int    `yaml:"log_autovacuum_min_duration"`
	AutovacuumMaxWorkers            int    `yaml:"autovacuum_max_workers"`
	AutovacuumNaptime               string `yaml:"autovacuum_naptime"`
	AutovacuumVacuumThreshold       int    `yaml:"autovacuum_vacuum_threshold"`
	AutovacuumAnalyzeThreshold      int    `yaml:"autovacuum_analyze_threshold"`
	AutovacuumVacuumScaleFactor     string `yaml:"autovacuum_vacuum_scale_factor"`
	AutovacuumAnalyzeScaleFactor    string `yaml:"autovacuum_analyze_scale_factor"`
	AutovacuumFreezeMaxAge          int    `yaml:"autovacuum_freeze_max_age"`
	AutovacuumMultixactFreezeMaxAge int    `yaml:"autovacuum_multixact_freeze_max_age"`
	AutovacuumVacuumCostDelay       string `yaml:"autovacuum_vacuum_cost_delay"`
	AutovacuumVacuumCostLimit       int    `yaml:"autovacuum_vacuum_cost_limit"`
}

type ConfigLine struct {
	OrigLine         string
	UseOrig          bool
	Key              string
	Value            string
	PostValueComment string
}
