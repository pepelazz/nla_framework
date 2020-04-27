<template>
  <div>
    <div v-if='!isRegisterSuccess' class="q-gutter-md">
      <q-input v-for="(fld, index) in flds" :key="fld.model" outlined :type='fld.type' :label="fld.label"
               v-model="regForm[fld.model]" :autofocus='index===0'>
        <template v-slot:prepend>
          <q-icon :name="fld.icon"/>
        </template>
      </q-input>
      <div class="row wrap justify-center items-start content-start q-gutter-md" style="margin-left: 0">
        <q-btn class="col" outline color="secondary" @click='$emit("cancel")'>отмена</q-btn>
        <q-btn class="col" color="primary" @click="login">ok</q-btn>
      </div>
    </div>
    <!--СООБЩЕНИЕ ПОСЛЕ ОТПРАВКИ ФОРМЫ РЕГИСТРАЦИИ -->
    <div v-if='isRegisterSuccess'>
      <div style="padding: 20px">Проверьте указанный в регистрации email.<br>Должно прийти письмо для
        подтверждения почтового адреса.
      </div>
    </div>
  </div>
</template>

<script>
    export default {
        data() {
            return {
                regForm: {},
                flds: [
                    {model: 'login', label: 'email', type: 'email', icon: 'email'},
                    {model: 'password', label: 'пароль', type: 'password', icon: 'lock'},
                    {model: 'passwordRepeat', label: 'повторите пароль', type: 'password', icon: 'lock'},
                    {model: 'last_name', label: 'фамилия', type: 'text', icon: 'person'},
                    {model: 'first_name', label: 'имя', type: 'text', icon: 'person_outline'},
                ],
                isRegisterSuccess: false,
            }
        },
        methods: {
            login() {
                // валидация полей формы регистрации
                // -- валидация email
                if (!validateEmail(this.regForm.login)) {
                    this.$q.notify({message: 'Поле "email" заполнено неверно', type: 'negative', position: 'top-right'})
                    return
                }
                // -- валидация пароля
                if (!this.regForm.password || this.regForm.password.length < 7) {
                    this.$q.notify({
                        message: 'Пароль должен быть больше 7 знаков',
                        type: 'negative',
                        position: 'top-right'
                    })
                    return
                }
                if (this.regForm.password !== this.regForm.passwordRepeat) {
                    this.$q.notify({
                        message: 'Повторно введенный пароль не совпадает с первым вариантом',
                        type: 'negative',
                        position: 'top-right'
                    })
                    return
                }
                // -- валидация имя (если поле указано в форме регистрации)
                if (!this.regForm.first_name) {
                    this.$q.notify({
                        message: 'Необходимо заполнить поле "имя"',
                        type: 'negative',
                        position: 'top-right'
                    })
                    return
                }
                // -- валидация фамилии (если поле указано в форме регистрации)
                if (!this.regForm.last_name) {
                    this.$q.notify({
                        message: 'Необходимо заполнить поле "фамилия"',
                        type: 'negative',
                        position: 'top-right'
                    })
                    return
                }
                // добавляем флаг, что это регистрация нового пользователя, а не авторизация по логину и паролю
                let params = Object.assign({is_register: true}, this.regForm)
                this.$utils.postApiRequest({url: '/auth/email', params, isShowError: false}).subscribe(res => {
                    if (res.ok) {
                        this.isRegisterSuccess = true
                    } else {
                        if (res.message.includes('email already exist')) {
                            this.$q.notify({
                                message: 'Пользователь с таким email уже зарегестрирован',
                                type: 'negative',
                                position: 'top-right'
                            })
                        } else {
                            this.$q.notify({message: res.message, type: 'negative', position: 'top-right'})
                        }
                    }
                })
            },
        },
        mounted() {
            this.flds.map(v => this.$set(this.regForm, v, null))
        }
    }
    const validateEmail = (email) => {
        let re = /\S+@\S+\.\S+/
        return re.test(email)
    }
</script>
