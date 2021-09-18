<template>
  <div class="q-gutter-md">
    <q-input outlined type='email' label="email" v-model="email" :autofocus=true>
      <template v-slot:prepend>
        <q-icon name="email"/>
      </template>
    </q-input>
    <q-input outlined type='password' :label="$t('auth.password')" v-model="password" @keydown.enter.prevent="login">
      <template v-slot:prepend>
        <q-icon name="lock"/>
      </template>
    </q-input>
    <div>
      <q-btn color="primary" class="full-width" @click='login' @keyup.enter='login'>{{$t('auth.login')}}</q-btn>
    </div>
    <div class="row wrap justify-center items-start content-start">
      <q-btn class="col" flat color="secondary" style="margin-bottom: 10px" @click='$emit("register")'>{{$t('auth.register')}}</q-btn>
      <q-btn class="col" flat color="secondary" style="margin-bottom: 10px" @click='$emit("passwordRecover")'>{{$t('auth.forgot_password')}}?</q-btn>
    </div>
  </div>
</template>

<script>
    export default {
        data() {
            return {
                email: null,
                password: null,
            }
        },
        methods: {
            login() {
                // валидация полей формы регистрации
                // -- валидация email
                if (!validateEmail(this.email)) {
                    this.$q.notify({message: this.$t('auth.invalid_email'), type: 'negative', position: 'top-right'})
                    return
                }
                // -- валидация пароля
                if (!this.password || this.password.length < 7) {
                    this.$q.notify({
                        message: this.$t('auth.invalid_password_must_be_more_7'),
                        type: 'negative',
                        position: 'top-right'
                    })
                    return
                }
                this.$utils.postApiRequest({
                    url: '/auth/email',
                    params: {login: this.email, password: this.password},
                    isShowError: false,
                }).subscribe(res => {
                    if (res.ok) {
                        this.$currentUser.login({user: res.result})
                    } else {
                        if (res.message.includes('user not found')) {
                            this.$q.notify({
                                message: this.$t('auth.user_not_found'),
                                type: 'negative',
                                position: 'top-right'
                            })
                        } else if (res.message.includes('wrong password')) {
                            this.$q.notify({message: this.$t('auth.wrong_password'), type: 'negative', position: 'top-right'})
                        } else {
                            this.$q.notify({message: res.message, type: 'negative', position: 'top-right'})
                        }
                    }
                })
            },
        }
    }
    const validateEmail = (email) => {
        let re = /\S+@\S+\.\S+/
        return re.test(email)
    }
</script>
