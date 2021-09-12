package types

func (p *ProjectType) FillI18n()  {
	p.I18n.LangList = []string{"en", "ru"}
	p.I18n.Data = map[string]map[string]map[string]string{
		"ru": {
			"message": {
				"ok": "ok",
				"save": "сохранить",
				"cancel": "отмена",
				"delete": "удалить",
				"edit": "редактировать",
				"send": "отправить",
				"file": "файл",
				"files": "файлы",
				"photo": "фото",
			},
			"profile": {
				"breadcrumb_label": "Редактирование профиля",
				"last_name": "фамилия",
				"first_name": "имя",
				"phone": "телефон",
				"avatar": "фото",
				"exit": "выйти",
			},
			"user": {
				"name_plural": "пользователи",
				"name_plural_deleted": "удаленные пользователи",
				"last_name": "фамилия",
				"first_name": "имя",
				"roles": "роли",
				"state": "статус",
				"phone": "телефон",
				"grade": "должность",
				"photo": "фото",
				"email": "email",
			},
			"menu": {
				"user": "пользователи",
			},
		},
		"en": {
			"message": {
				"ok": "ok",
				"save": "save",
				"cancel": "cancel",
				"delete": "delete",
				"edit": "edit",
				"send": "send",
				"file": "file",
				"files": "files",
				"photo": "photo",
			},
			"profile": {
				"breadcrumb_label": "Profile edit",
				"last_name": "last name",
				"first_name": "first name",
				"phone": "phone",
				"avatar": "photo",
				"exit": "exit",
			},
			"user": {
				"name_plural": "users",
				"name_plural_deleted": "deleted users",
				"last_name": "last name",
				"first_name": "first name",
				"roles": "roles",
				"state": "state",
				"phone": "phone",
				"grade": "grade",
				"photo": "photo",
				"email": "email",
			},
			"menu": {
				"user": "users",
			},
		},
	}
}

func (p *ProjectType) AddI18nLanguage(lang string) {
	p.I18n.LangList = append(p.I18n.LangList, lang)
}
