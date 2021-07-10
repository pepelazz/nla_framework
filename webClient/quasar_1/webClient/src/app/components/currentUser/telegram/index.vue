<template>
  <div>
    <vue-telegram-login
      mode="callback"
      :telegram-login="$config.telegram.botName" size="medium"
      @callback="cb"
    />
  </div>
</template>

<script>
    import vueTelegramLogin from './vueTelegramLogin'

    export default {
        props: ['isRegistered'],
        components: {vueTelegramLogin},
        data() {
            return {
                isDone: false,
            }
        },
        methods: {
            cb(user) {
                if (this.isRegistered || this.isDone) return
                this.$utils.postApiRequest({url: "/api/telegram_auth", params: JSON.stringify(user)}).subscribe(res => {
                    if (res.ok) {
                        this.isDone = false
                        this.$q.notify({
                            message: 'аккаунт в телеграм успешно зарегестрирован',
                            type: 'positive',
                            position: 'bottom-right'
                        })
                    }
                })
            }
        }
    }
</script>
