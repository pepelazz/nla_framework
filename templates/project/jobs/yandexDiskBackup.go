package jobs

import (
	"[[.Config.LocalProjectPath]]/yandexDiskBackup"
	"time"
)

func startYandexDiskBackup()  {
	go func() {
		yandexDiskBackup.CreateBackup()
		time.Sleep([[.Config.Backup.ToYandexDisk.Period]]* time.Minute)
	}()
}
