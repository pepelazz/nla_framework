<template>
  <div class="layout">

    <!--поля для восстановления пароля-->
    <div v-if="isShowForm && !isSuccess" class="row justify-center q-ma-lg">
      <q-card class="col-md-4 col-xs-12 q-pa-lg" style="min-width: 35%">
        <q-toolbar>
          <q-toolbar-title><span class="text-weight-bold" style="text-transform: capitalize">{{$t('auth.password_recovery')}}</span></q-toolbar-title>
        </q-toolbar>
        <div class="q-gutter-md q-mt-md">
          <q-input v-for="fld in flds" :key="fld.model" outlined :type='fld.type' :label="fld.label"
                   v-model="regForm[fld.model]">
            <template v-slot:prepend>
              <q-icon :name="fld.icon"/>
            </template>
          </q-input>
          <div class="row wrap justify-center items-start content-start q-gutter-md" style="margin-left: 0">
            <q-btn class="col" color="primary" @click="changePassword">ok</q-btn>
          </div>
        </div>
      </q-card>
    </div>

    <!--сообщение когда пароль успешно изменен-->
    <q-banner v-if="isSuccess" class="bg-primary text-white q-ma-lg">
      <template v-slot:avatar>
        <q-icon name="lock" color="white"/>
      </template>
      {{$t('auth.password_recovery_success_message')}}
      <template v-slot:action>
        <q-btn flat color="white" :label="$t('auth.password_recovery_go_to_login')" @click="$router.push(homeUrl)"/>
      </template>
    </q-banner>

    <!--сообщения в случае невалидного токена-->
    <q-banner v-if="isTokenNotValid" class="bg-red text-white q-ma-lg">
      <template v-slot:avatar>
        <q-icon name="error" color="white"/>
      </template>
      Ссылка уже неактивна. Возможные причины:<br><br>
      Вы уже успешно поменяли пароль<br>
      Истекло время в течении, которого можно было воспользоваться этой ссылкой
      <template v-slot:action>
        <q-btn flat color="white" label="Перейти к авторизации" @click="$router.push(homeUrl)"/>
      </template>
    </q-banner>

  </div>
</template>

<script>
    export default {
        computed: {
            homeUrl() {
                return `/`
            }
        },
        data() {
            return {
                regForm: {},
                flds: [
                    {model: 'password', label: this.$t('auth.password'), type: 'password', icon: 'lock'},
                    {model: 'passwordRepeat', label: this.$t('auth.password_repeat'), type: 'password', icon: 'lock'},
                ],
                isTokenNotValid: null,
                isShowForm: false,
                isSuccess: false,
                token: null,
            }
        },
        methods: {
            changePassword() {
                // -- валидация пароля
                if (!this.regForm.password || this.regForm.password.length < 7) {
                    this.$q.notify({
                        message: this.$t('auth.invalid_password_must_be_more_7'),
                        type: 'negative',
                        position: 'top-right'
                    })
                    return
                }
                if (this.regForm.password !== this.regForm.passwordRepeat) {
                    this.$q.notify({
                        message: this.$t('auth.invalid_password_wrong_repeat'),
                        type: 'negative',
                        position: 'top-right'
                    })
                    return
                }
                this.$utils.postApiRequest({
                    url: '/auth/email_auth_recover_password',
                    params: {token: this.token, password: this.regForm.password}
                }).subscribe(res => {
                    res.ok ? this.isSuccess = true : this.$q.notify({
                        message: res.message,
                        type: 'negative',
                        position: 'top-right'
                    })
                })
            }
        },
        created() {
            // извлекаем токен из url
            let url = window.location.search
            if (url.split('?').length > 1) {
                let params = url.split('?')[1].split('=')
                if (params[0] === 't' && params.length > 1) {
                    this.token = params[1]
                    // проверяем токен в базе
                    this.$utils.postApiRequest({
                        url: '/auth/email_auth_recover_password',
                        params: {token: this.token, is_token_check: true},
                        isShowError: false
                    }).subscribe(res => {
                        res.ok ? this.isShowForm = true : this.isTokenNotValid = true
                    })
                }
            }
        }
    }
</script>
