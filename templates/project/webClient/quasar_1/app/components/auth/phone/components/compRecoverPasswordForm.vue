<template>
  <div>
    <div v-if="processState == 'step1'" class="q-gutter-md" @keydown.enter.prevent="ok">
      <q-input outlined type='text' label="Ваш телефон" mask="+# (###) ### - ####" v-model="phone" :autofocus=true>
        <template v-slot:prepend>
          <q-icon name="phone"/>
        </template>
      </q-input>
      <div class="row wrap justify-center items-start content-start q-gutter-md" style="margin-left: 0">
        <q-btn class="col" outline color="secondary" @click='$emit("cancel")'>{{$t('message.cancel')}}</q-btn>
        <q-btn class="col" color="primary" @click="ok">ok</q-btn>
      </div>
    </div>
    <!--СООБЩЕНИЕ ДЛЯ ВВОДА ПРОВЕРОЧНОГО КОДА, ПРИСЛАННОГО ПО SMS-->
    <div v-if="processState == 'step2'" class="q-gutter-md" @keydown.enter.prevent="changePassword">
      <q-input outlined type='text' label="введите код из sms"
               v-model="smsCode" autofocus mask="######">
        <template v-slot:prepend>
          <q-icon name="fas fa-sms"/>
        </template>
      </q-input>
      <q-input outlined type='password' label="новый пароль"
               v-model="password">
        <template v-slot:prepend><q-icon name="lock"/></template>
      </q-input>
      <q-input outlined type='password' label="повторите пароль"
               v-model="passwordRepeat">
        <template v-slot:prepend><q-icon name="lock"/></template>
      </q-input>
      <div class="row wrap justify-center items-start content-start q-gutter-md" style="margin-left: 0">
        <q-btn class="col" outline color="secondary" @click='$emit("cancel")'>{{$t('message.cancel')}}</q-btn>
        <q-btn class="col" color="primary" @click="changePassword">ok</q-btn>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        phone: null,
        processState: 'step1',
        smsCode: null,
        password: null,
        passwordRepeat: null,
      }
    },
    methods: {
      ok() {
        // -- валидация
        if (!this.phone) {
          this.$q.notify({message: 'Поле "номер телефона" заполнено неверно', type: 'negative', position: 'top-right'})
          return
        }
        // оставляем только цифры
        this.phone = this.phone.replace(/\D/g, '')
        // если первая цифра 8, то заменяем на 7
        if (this.phone.charAt(0) === '8') this.phone = '7' + this.phone.slice(1)
        if (this.phone.length !== 11) {
          this.$q.notify({message: 'Поле "номер телефона" заполнено неверно', type: 'negative', position: 'top-right'})
          return
        }
        this.$utils.postApiRequest({
          url: '/auth/phone_auth_start_recover_password',
          params: {phone: this.phone},
          isShowError: false
        }).subscribe(res => {
          if (res.ok) {
            this.processState = 'step2'
          } else {
            if (res.message.includes('user not found')) {
              this.$q.notify({
                message: 'Пользователь с таким телефоном не зарегестрирован',
                type: 'negative',
                position: 'top-right'
              })
            } else {
              this.$q.notify({message: res.message, type: 'negative', position: 'top-right'})
            }
          }
        })
      },
      changePassword() {
        if (!this.smsCode) {
          this.$q.notify({
            message: 'Необходимо заполнить поле с кодом из sms',
            type: 'negative',
            position: 'top-right'
          })
          return
        }
        // -- валидация пароля
        if (!this.password || this.password.length < 7) {
          this.$q.notify({
            message: 'Пароль должен быть больше 7 знаков',
            type: 'negative',
            position: 'top-right'
          })
          return
        }
        if (this.password !== this.passwordRepeat) {
          this.$q.notify({
            message: 'Повторно введенный пароль не совпадает с первым вариантом',
            type: 'negative',
            position: 'top-right'
          })
          return
        }
        this.$utils.postApiRequest({
          url: '/auth/phone_auth_recover_password',
          params: {phone: this.phone.replace(/\D/g, ''), token: this.smsCode, password: this.password},
          isShowError: false
        }).subscribe(res => {
          if (res.ok) {
            this.$emit("cancel")
            this.$q.notify({
              message: 'Пароль успешно изменен',
              type: 'positive',
              position: 'bottom-right'
            })
          } else {
            if (res.message.includes('user not found')) {
              this.$q.notify({
                message: 'Пользователь с таким телефоном не зарегестрирован',
                type: 'negative',
                position: 'top-right'
              })
            }
            else if (res.message.includes('invalid token')) {
              this.$q.notify({
                message: 'неверный код из sms',
                type: 'negative',
                position: 'top-right'
              })
            } else {
              this.$q.notify({message: res.message, type: 'negative', position: 'top-right'})
            }
          }
        })
      },
    }
  }
</script>
