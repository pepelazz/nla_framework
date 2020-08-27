package main

import (
	"fmt"
	"time"
)

var ynxUrl = "https://cloud-api.yandex.net/v1/disk"
var authToken = "[[.Config.Backup.ToYandexDisk.Token]]"

func main()  {

		for {
			createBackup()
			time.Sleep([[.Config.Backup.ToYandexDisk.Period]]* time.Minute)
	}
}

func createBackup()  {
	// создаем папку на яндексе для проекта
	err := createFolder("[[.Config.Backup.ToYandexDisk.Path]]")
	if err != nil {
		fmt.Printf("error createFolder [[.Config.Backup.ToYandexDisk.Path]]")
	}
	// получаем адрес файла с бэкапом
	fileName, err := getBackupFile()
	if err != nil {
		fmt.Printf("error getBackupFile %s err:%s", fileName, err)
		return
	}
	// копируем файл на сервер
	err = uploadFile(fileName, "[[.Config.Backup.ToYandexDisk.Path]]/"+fileName)
	if err != nil {
		fmt.Printf("error uploadFile %s error:%s", fileName, err)
		return
	}

	// удаляем файл на сервере с таймаутом
	time.Sleep(1 * time.Minute)
	removeBackupFile(fileName)

	// удаляем старые файлы на яндекс диске
	err = removeOldBackupsOnServer()
	if err != nil {
		fmt.Printf("error removeBackupFile error:%s", err)
	}

}