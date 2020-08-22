<template>
  <div>
    <div v-if="!isSuccess" class="q-gutter-md">
      <q-input outlined type='text' label="Ваш телефон" mask="+# (###) ### - ####" v-model="phone" :autofocus=true>
        <template v-slot:prepend>
          <q-icon name="phone"/>
        </template>
      </q-input>
      <div class="row wrap justify-center items-start content-start q-gutter-md" style="margin-left: 0">
        <q-btn class="col" outline color="secondary" @click='$emit("cancel")'>отмена</q-btn>
        <q-btn class="col" color="primary" @click="ok">ok</q-btn>
      </div>
    </div>
    <!--СООБЩЕНИЕ ДЛЯ ВВОДА ПРОВЕРОЧНОГО КОДА, ПРИСЛАННОГО ПО SMS-->
    <div v-if='isSuccess' class="q-gutter-md">
      <q-input outlined type='text' label="введите код из sms"
               v-model="smsCode" autofocus mask="######">
        <template v-slot:prepend>
          <q-icon name="fas fa-sms"/>
        </template>
      </q-input>
      <div class="row wrap justify-center items-start content-start q-gutter-md" style="margin-left: 0">
        <q-btn class="col" outline color="secondary" @click='$emit("cancel")'>отмена</q-btn>
        <q-btn class="col" color="primary" @click="checkCode">ok</q-btn>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        phone: null,
        isSuccess: false,
        smsCode: null,
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
            this.isSuccess = true
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
      checkCode() {}
    }
  }
</script>
