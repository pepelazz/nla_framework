package templates

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/templates/tmplGenerateStep2"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"strings"
	"text/template"
)

func WriteProjectFiles(p types.ProjectType, tmplMap map[string]*template.Template)  {
	for name, t := range tmplMap {
		if strings.HasPrefix(name, "project_") {
			filename := strings.TrimPrefix(name, "project_")
			path := ".."
			if filename == "config.toml" || filename == "main.go" {
				path = "../src"
			}
			err := ExecuteToFile(t, p, path, filename)
			utils.CheckErr(err, fmt.Sprintf("'project' ExecuteToFile '%s'", name))
		}
	}

	// генерим шаблоны, которые указаны дополнительно на уровне проекта. Без относительно конкретных документов
	for _, m := range p.Sql.Methods {
		for _, v := range m {
			if len(v.Tmpl.Source) > 0 && len(v.Tmpl.Dist) > 0 {
				distPath, filename := utils.PathExtractFilename(v.Tmpl.Dist)
				distPath = "../src" + distPath
				t, err := template.New(filename).Delims("[[", "]]").ParseFiles(v.Tmpl.Source)
				utils.CheckErr(err, "p.Sql.Methods")

				err = ExecuteToFile(t, p, distPath, filename)
				utils.CheckErr(err, fmt.Sprintf("'project' ExecuteToFile '%s'", filename))
			}
		}
	}

	projectTmplPath := "../../../pepelazz/projectGenerator/templates/project"
	ReadTmplAndPrint(p, projectTmplPath + "/types/main.go", "/types",  "main.go", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/types/config.go", "/types",  "config.go", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webServer/main.go", "/webServer",  "main.go", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webServer/apiCallPgFunc.go", "/webServer",  "apiCallPgFunc.go", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/sql/initialData.sql", "/sql/template/function/",  "initialData.sql", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/sql/user_trigger_after.sql", "/sql/template/function/_User/",  "user_trigger_after.sql", template.FuncMap{"PrintUserAfterTriggerUpdateLinkedRecords": types.PrintUserAfterTriggerUpdateLinkedRecords})
	ReadTmplAndPrint(p, projectTmplPath + "/sql/01_User/main.toml", "/sql/model/01_User",  "main.toml", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/sql/03_UserTempEmailAuth/main.toml", "/sql/model/03_UserTempEmailAuth",  "main.toml", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/jobs/main.go", "/jobs",  "main.go", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/pg/pgListener.go", "/pg",  "pgListener.go", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webClient/index.template.html", "/webClient/src",  "index.template.html", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webClient/App.vue", "/webClient/src",  "App.vue", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/users/roles.js", "/webClient/src/app/components/users",  "roles.js", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/users/item.vue", "/webClient/src/app/components/users",  "item.vue", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/users/index.vue", "/webClient/src/app/components/users",  "index.vue", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/currentUser/profile.vue", "/webClient/src/app/components/currentUser",  "profile.vue", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/currentUser/messages/list.vue", "/webClient/src/app/components/currentUser/messages",  "list.vue", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/home.vue", "/webClient/src/app/components",  "home.vue", nil)
	ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/auth/loginPage.vue", "/webClient/src/app/components/auth",  "loginPage.vue", nil)

	if p.Config.Auth.ByPhone {
		ReadTmplAndPrint(p, projectTmplPath + "/sql/01_User/user_get_by_phone_with_password.sql", "/sql/template/function/_User",  "user_get_by_phone_with_password.sql", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/sql/03_UserTempEmailAuth/user_temp_phone_auth_create.sql", "/sql/template/function/_UserTempEmailAuth",  "user_temp_phone_auth_create.sql", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/sql/03_UserTempEmailAuth/user_temp_phone_auth_check_sms_code.sql", "/sql/template/function/_UserTempEmailAuth",  "user_temp_phone_auth_check_sms_code.sql", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/webServer/auth/phone.go", "/webServer/auth",  "phone.go", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/auth/phone/phoneAuthBtn.vue", "/webClient/src/app/components/auth/phone",  "phoneAuthBtn.vue", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/auth/phone/components/compLoginForm.vue", "/webClient/src/app/components/auth/phone/components",  "compLoginForm.vue", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/auth/phone/components/compRecoverPasswordForm.vue", "/webClient/src/app/components/auth/phone/components",  "compRecoverPasswordForm.vue", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/webClient/app/components/auth/phone/components/compRegisterForm.vue", "/webClient/src/app/components/auth/phone/components",  "compRegisterForm.vue", nil)
	}

	if p.IsTelegramIntegration() {
		ReadTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/telegram/telegramAuth.go", "/webServer", "telegramAuth.go", nil)
		ReadTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/telegram/user_telegram_auth.sql", "/sql/template/function/_User", "user_telegram_auth.sql", nil)
		ReadTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/telegram/user_get_by_telegram_id.sql", "/sql/template/function/_User", "user_get_by_telegram_id.sql", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/tgBot/main.go", "/tgBot", "main.go", nil)
	}

	if p.IsBackupOnYandexDisk() {
		ReadTmplAndPrint(p, projectTmplPath + "/yandexDiskBackup/main.go", "/yandexDiskBackup",  "main.go", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/yandexDiskBackup/yandexApi.go", "/yandexDiskBackup",  "yandexApi.go", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/yandexDiskBackup/dbBackup.go", "/yandexDiskBackup",  "dbBackup.go", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/yandexDiskBackup/systemdService.service", "/yandexDiskBackup",  p.Config.Postgres.DbName + "_yandexBackup.service", nil)
		ReadTmplAndPrint(p, projectTmplPath + "/yandexDiskBackup/startYandexBackupService.sh", "/yandexDiskBackup","startYandexBackupService.sh", nil)
	}

	// в случае коннекта к Битрикс генерим файлы
	if p.IsBitrixIntegration() {
		ReadTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/bitrix/bitrixMain.go", "/bitrix", "main.go", nil)
		//sourcePath := "../../../pepelazz/projectGenerator/templates/integrations/bitrix/bitrixMain.go"
		//t, err := template.New("bitrixMain.go").Funcs(funcMap).Delims("[[", "]]").ParseFiles(sourcePath)
		//utils.CheckErr(err, "bitrixMain.go")
		//distPath := fmt.Sprintf("%s/bitrix", p.DistPath)
		//err = ExecuteToFile(t, p, distPath, "main.go")
		//utils.CheckErr(err, fmt.Sprintf("'project' ExecuteToFile '%s'", "bitrix/main.go"))
	}

	// в случае коннекта к 1 Odata генерим файлы
	if p.IsOdataIntegration() {
		ReadTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/odata/main.go", "/odata", "main.go", nil)
		ReadTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/odata/odataQueryType.go", "/odata", "odataQueryType.go", nil)
	}
}

func OtherTemplatesGenerate(p types.ProjectType)  {
	tmplGenerateStep2.TasksTmpl(p)
	// добавляем функции в plugin/utils.js
	tmplGenerateStep2.PluginUtilsJs(p)
}

func ReadTmplAndPrint(p types.ProjectType, sourcePath, distPath, filename string, addFuncMap template.FuncMap) {
	fMap := template.FuncMap{}
	for k, v := range funcMap {
		fMap[k] = v
	}
	for k, v := range addFuncMap {
		fMap[k] = v
	}
	// проверяем возможность того, что sourcePath был переопределен
	if newSourcePath, ok := p.OverridePathForTemplates[fmt.Sprintf("%s/%s", distPath, filename)]; ok {
		sourcePath = newSourcePath
	}
	_, sourceFilename := utils.PathExtractFilename(sourcePath)
	t, err := template.New(sourceFilename).Funcs(fMap).Delims("[[", "]]").ParseFiles(sourcePath)
	utils.CheckErr(err, "readFileWithDist")
	err = ExecuteToFile(t, p, p.DistPath + distPath, filename)
	utils.CheckErr(err, fmt.Sprintf("ReadTmplAndPrint ExecuteToFile '%s/%s'", distPath, filename))
}

