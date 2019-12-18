package logging

import (
	"fmt"
	"time"

	"vcard/pkg/setting"
)

// getLogFilePath get the log file save path
func GetLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

// getLogFileName get the save name of the log file
func GetLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}
