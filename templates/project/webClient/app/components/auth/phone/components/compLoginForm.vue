<template>
  <div class="q-gutter-md">
    <q-input outlined mask="(###) ### - ####" label="номер телефона" v-model="phone" :autofocus=true>
      <template v-slot:prepend>
        <q-icon name="phone"/>
      </template>
    </q-input>
    <q-input outlined type='password' label="пароль" v-model="password">
      <template v-slot:prepend>
        <q-icon name="lock"/>
      </template>
    </q-input>
    <div>
      <q-btn color="primary" class="full-width" @click='login' @keyup.enter='login'>войти</q-btn>
    </div>
    <div class="row wrap justify-center items-start content-start">
      <q-btn class="col" flat color="secondary" style="margin-bottom: 10px" @click='$emit("register")'>
        зарегистрироваться
      </q-btn>
      <q-btn class="col" flat color="secondary" style="margin-bottom: 10px" @click='$emit("passwordRecover")'>забыли
        пароль?
      </q-btn>
    </div>
  </div>
</template>

<script>
    export default {
        data() {
            return {
                phone: null,
                password: null,
            }
        },
        methods: {
            login() {
                // валидация полей формы регистрации
                // -- валидация телефона
                if (!this.phone) {
                  this.$q.notify({message: 'Поле "номер телефона" заполнено неверно', type: 'negative', position: 'top-right'})
                  return
                }
                // оставляем только цифры
                this.phone = this.phone.replace(/\D/g, '')
                if (this.phone.length !== 11) {
                  this.$q.notify({message: 'Поле "номер телефона" заполнено неверно', type: 'negative', position: 'top-right'})
                  return
                }
                // если первая цифра 8, то заменяем на 7
                if (this.phone.charAt(0) === '8') this.phone = '7' + this.phone.slice(1)
                // -- валидация пароля
                if (!this.password || this.password.length < 7) {
                    this.$q.notify({
                        message: 'Пароль должен быть больше 7 знаков',
                        type: 'negative',
                        position: 'top-right'
                    })
                    return
                }
                this.$utils.postApiRequest({
                    url: '/auth/phone',
                    params: {login: this.phone, password: this.password},
                    isShowError: false,
                }).subscribe(res => {
                    if (res.ok) {
                        this.$currentUser.login({user: res.result})
                    } else {
                        if (res.message.includes('user not found')) {
                            this.$q.notify({
                                message: 'Пользователь с таким номером телефона не найден',
                                type: 'negative',
                                position: 'top-right'
                            })
                        } else if (res.message.includes('wrong password')) {
                            this.$q.notify({message: 'Неверный пароль', type: 'negative', position: 'top-right'})
                        } else {
                            this.$q.notify({message: res.message, type: 'negative', position: 'top-right'})
                        }
                    }
                })
            },
        }
    }
</script>
