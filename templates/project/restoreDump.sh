#!/bin/bash

# функция выхода из скрипта при ошибке
is_err () {
    [ $? -ne 0 ]
}

# функция выхода из скрипта при ошибке
is_err () {
    [ $? -ne 0 ]
}

echo -e "\033[0;32m STEP1: create database dump...\033[0m"
ssh {{.Config.WebServer.Username}}@{{.Config.WebServer.Ip}} << EOF
    cd {{.Config.WebServer.Path}}
    docker exec -t {{.Name | ToLower}}_postgres_1 pg_dumpall -c -U postgres  > {{.Name | ToLower}}_dump
EOF
if is_err; then return; fi

echo -e "\033[0;32m STEP2: copy file from server...\033[0m"
scp {{.Config.WebServer.Username}}@{{.Config.WebServer.Ip}}:/{{.Config.WebServer.Path}}/{{.Name | ToLower}}_dump .

# запускаем докер
docker-compose --file docker-compose.dev.yml up -d

# удаляем базу
echo -e "\033[0;32m STEP1: delete database...\033[0m"
sleep 5
docker exec -t {{.Name | ToLower}}_postgres_1 psql -U postgres -c 'DROP DATABASE {{.Config.Postgres.DbName}}'

# восстанавливаем базу
echo -e "\033[0;32m STEP2: restore database...\033[0m"
sleep 5
cat {{.Name | ToLower}}_dump | docker exec -i {{.Name | ToLower}}_postgres_1 psql -U postgres

# останавливаем докер
docker-compose stop