#!/bin/bash
chmod a+rwx -R [[.Config.WebServer.Path]]/yandexDiskBackup
cp [[.Config.WebServer.Path]]/src/yandexDiskBackup/[[.Config.Postgres.DbName]]_yandexBackup.service /etc/systemd/system/[[.Config.Postgres.DbName]]_yandexBackup.service
chmod a+rwx -R /etc/systemd/system/[[.Config.Postgres.DbName]]_yandexBackup.service
systemctl daemon-reload
systemctl enable [[.Config.Postgres.DbName]]_yandexBackup.service
systemctl start [[.Config.Postgres.DbName]]_yandexBackup