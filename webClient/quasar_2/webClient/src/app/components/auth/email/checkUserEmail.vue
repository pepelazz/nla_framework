<template>
  <div class="layout">

    <!--СООБЩЕНИЯ В СЛУЧАЕ НЕВАЛИДНОГО ТОКЕНА-->
    <div style="margin-top: 20px" v-if="isTokenNotValid">
      <div class="row justify-center">
        <img src="https://a.radikal.ru/a01/1806/fc/49f02b5a592a.jpg" alt="">
      </div>
      <div class="row justify-center">
        <q-banner inline-actions class="bg-grey-3 q-pa-md" v-html="$t('auth.check_user_email_message')"/>
      </div>
      <div class="row justify-center" style="margin-top: 50px">
        <q-btn @click="$router.push(homeUrl)" color="primary">Перейти к авторизации</q-btn>
      </div>
    </div>
  </div>
</template>

<script>
    export default {
        props: ['isCheckUserEmail'],
        computed: {
            homeUrl() {
                return `/`
            }
        },
        data() {
            return {
                isTokenNotValid: null
            }
        },
        mounted() {
            // извлекаем токен из url
            let url = window.location.search
            if (url.split('?').length > 1) {
                let params = url.split('?')[1].split('=')
                if (params[0] === 't' && params.length > 1) {
                    // проверяем токен в базе
                    this.$utils.postApiRequest({
                        url: `/auth/check_user_email`,
                        params: {token: params[1]},
                        isShowError: false,
                    }).subscribe(res => {
                        if (res.ok) {
                            this.$currentUser.login({user: res.result})
                            this.$router.push(this.homeUrl)
                        } else {
                            this.isTokenNotValid = true
                        }
                    })
                }
            }
        }
    }
</script>
