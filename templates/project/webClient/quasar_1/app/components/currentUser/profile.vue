<template>
  <q-page padding>

    <comp-breadcrumb :list="[{label:'Редактирование профиля'}]"/>

    <div v-if="item" class="q-mt-sm">
      <!--  поля формы    -->
      <div class="row q-col-gutter-md q-mb-sm" v-for="fldRow in flds">
        <comp-fld v-for="fld in fldRow" :key='fld.name'
                  :fld="item[fld.name]"
                  :type="fld.type"
                  @update="item[fld.name] = $event"
                  :label="fld.label"
                  :selectOptions="fld.selectOptions ? fld.selectOptions() : []"
                  :ajaxSelectTitle="item[fld.ajaxSelectTitle]"
                  :columnClass="fld.columnClass"
                  :compName="fld.compName"
                  :ext="fld.ext"
        />
      </div>
      <telegram-login v-if="$config.telegram && $config.telegram.botName" :isRegistered="currentUser.options.telegram_id"/>
      <!--  кнопки   -->
      <comp-item-btn-save @save="save" @cancel="$router.push(docUrl)"/>
    </div>
  </q-page>
</template>

<script>
    import telegramLogin from './telegram/index'
    import currentUserMixin from '../../../app/mixins/currentUser'
    export default {
        mixins: [currentUserMixin],
        components: {telegramLogin},
        computed: {
            docUrl: () => '/',
        },
        data() {
            return {
                item: null,
               [[if .Vue.Hooks.Profile.Flds ]]
                  [[.Vue.Hooks.Profile.Flds]]
               [[- else]]
                flds: [
                    [
                        {name: 'last_name', type: 'string', label: 'Фамилия', required: true},
                        {name: 'first_name', type: 'string', label: 'Имя', required: true},
                    ],
                    [
                        {name: 'phone', type: 'phone', label: 'Телефон'},
                    ],
                    [[if not .Vue.IsHideUserAvatarUploader]][
                        {name: 'avatar', compName: 'comp-fld-img', label: 'Фото', ext: {fldName: 'avatar', uploadUrl: 'upload_profile_image', methodUpdate: 'current_user_update'}, columnClass: 'col-xs-6 col-sm-6 col-md-2'},
                    ],[[end]]
                ],
                optionsFlds: [],
               [[- end]]
            }
        },
        methods: {
            save() {
                this.$utils.saveItem.call(this, {
                    method: 'current_user_update',
                    itemForSaveMod: {
                        options: Object.assign(this.item.options || {}),
                    },
                    resultModify: (res) => {
                        // для обновления currentUser выполняем операцию login
                        this.login()
                        if (this.optionsFlds) this.optionsFlds.map(fldName => res[fldName] = res.options[fldName])
                        return res
                    }
                })
            },
        },
        mounted() {
            this.item = this.currentUser
            if (this.optionsFlds) this.optionsFlds.map(fldName => this.$set(this.item, fldName, this.item.options[fldName]))
        }
    }
</script>
