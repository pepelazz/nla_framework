<template>
  <div>
    <div v-if="!isSuccess" class="q-gutter-md">
      <q-input outlined type='email' label="email" v-model="email" :autofocus=true>
        <template v-slot:prepend>
          <q-icon name="email"/>
        </template>
      </q-input>
      <div class="row wrap justify-center items-start content-start q-gutter-md" style="margin-left: 0">
        <q-btn class="col" outline color="secondary" @click='$emit("cancel")'>{{$t('message.cancel')}}</q-btn>
        <q-btn class="col" color="primary" @click="ok">ok</q-btn>
      </div>
    </div>
    <!--СООБЩЕНИЕ ПОСЛЕ ОТПРАВКИ СООБЩЕНИЯ О ВОССТАНОВЛЕНИИ ПАРОЛЯ -->
    <q-banner v-if='isSuccess' inline-actions class="bg-grey-3 q-pa-md">
      <div style="padding: 20px">{{$t('auth.password_recovery_message')}}</div>
    </q-banner>
  </div>
</template>

<script>
    export default {
        data() {
            return {
                email: null,
                isSuccess: false,
            }
        },
        methods: {
            ok() {
                // -- валидация email
                if (!validateEmail(this.email)) {
                    this.$q.notify({message: this.$t('invalid_email'), type: 'negative', position: 'top-right'})
                    return
                }
                this.$utils.postApiRequest({
                    url: '/auth/email_auth_start_recover_password',
                    params: {email: this.email},
                    isShowError: false
                }).subscribe(res => {
                    if (res.ok) {
                        this.isSuccess = true
                    } else {
                        if (res.message.includes('user not found')) {
                            this.$q.notify({
                                message: this.$t('auth.user_not_found'),
                                type: 'negative',
                                position: 'top-right'
                            })
                        } else {
                            this.$q.notify({message: res.message, type: 'negative', position: 'top-right'})
                        }
                    }
                })
            }
        }
    }
    const validateEmail = (email) => {
        let re = /\S+@\S+\.\S+/
        return re.test(email)
    }
</script>
