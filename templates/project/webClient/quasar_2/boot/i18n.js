import { createI18n } from 'vue-i18n'
import messages from '../i18n'

const i18n = createI18n({
    locale: [[- if .I18n.DefaultLang]] '[[.I18n.DefaultLang]]' [[else]] 'ru-RU'[[- end]],
    messages
})

export default ({ app }) => {
    // Set i18n instance on app
    app.use(i18n)
}

export { i18n }
