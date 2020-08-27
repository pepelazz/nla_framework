package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getBackupFile() (string, error) {

	fileName := fmt.Sprintf("[[.Config.Postgres.DbName]]_dump_%s.zip", time.Now().Format("2006_01_02_15_04"))
	// делаем бэкап базы на сервере
	cmd := exec.Command("sh", "-c", strings.Join([]string{
		//fmt.Sprintf("pg_dumpall -c -U postgres  > db_dump"),
		fmt.Sprintf("docker exec -t [[.Config.Postgres.DbName]]_postgres_1 pg_dumpall -c -U postgres  > db_dump"),
		fmt.Sprintf("zip %s db_dump", fileName), // архивируем бэкап
		fmt.Sprintf("rm db_dump"),               // удаляем бэкап
	}, ";"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return "", errors.New(fmt.Sprintf("file %s not created", fileName))
	}
	return fileName, nil
}

func removeBackupFile(path string) {
	cmd := exec.Command("sh", "-c", strings.Join([]string{
		fmt.Sprintf("rm %s", path),
	}, ";"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("removeBackupFile err %s\n", err)
	}
}

// удаляем файлы кроме указанного количества последних файлов
func removeOldBackupsOnServer() error {
	res, err := getResource("[[.Config.Backup.ToYandexDisk.Path]]")
	if err != nil {
		return err
	}
	for i, v := range res.Embedded.Items {
		//fmt.Printf("%v %s %s\n", i, v.Name, v.Path)
		if i >= [[.Config.Backup.ToYandexDisk.FilesCount]] {
			err = deleteFile(v.Path)
			if err != nil {
				fmt.Printf("deleteFile err: %s\n", err)
			}
		}
	}
	return nil
}
