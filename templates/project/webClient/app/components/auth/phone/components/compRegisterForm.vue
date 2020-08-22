<template>
    <div>
        <div v-if='!isRegisterSuccess' class="q-gutter-md">
            <q-input v-for="(fld, index) in flds" :key="fld.model" outlined :type='fld.type' :label="fld.label"
                     v-model="regForm[fld.model]" :autofocus='index===0' :mask="fld.mask">
                <template v-slot:prepend>
                    <q-icon :name="fld.icon"/>
                </template>
            </q-input>
            <div class="row wrap justify-center items-start content-start q-gutter-md" style="margin-left: 0">
                <q-btn class="col" outline color="secondary" @click='$emit("cancel")'>отмена</q-btn>
                <q-btn class="col" color="primary" @click="login">ok</q-btn>
            </div>
        </div>
        <!--СООБЩЕНИЕ ДЛЯ ВВОДА ПРОВЕРОЧНОГО КОДА, ПРИСЛАННОГО ПО SMS-->
        <div v-if='isRegisterSuccess' class="q-gutter-md">
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
                regForm: {},
                flds: [
                    {model: 'login', label: 'номер телефона', type: 'text', icon: 'phone', mask: '+# (###) ### - ####'},
                    {model: 'password', label: 'пароль', type: 'password', icon: 'lock'},
                    {model: 'passwordRepeat', label: 'повторите пароль', type: 'password', icon: 'lock'},
                    {model: 'last_name', label: 'фамилия', type: 'text', icon: 'person'},
                    {model: 'first_name', label: 'имя', type: 'text', icon: 'person_outline'},
                ],
                smsCode: null,
                isRegisterSuccess: false,
            }
        },
        methods: {
            login() {
                // валидация полей формы регистрации
                // -- валидация телефона
                if (!this.regForm.login) {
                    this.$q.notify({message: 'Поле "номер телефона" заполнено неверно', type: 'negative', position: 'top-right'})
                    return
                }
                // оставляем только цифры
                this.regForm.login = this.regForm.login.replace(/\D/g, '')
                if (this.regForm.login.length !== 11) {
                    this.$q.notify({message: 'Поле "номер телефона" заполнено неверно', type: 'negative', position: 'top-right'})
                    return
                }
                // если первая цифра 8, то заменяем на 7
                if (this.regForm.login.charAt(0) === '8') this.regForm.login = '7' + this.regForm.login.slice(1)
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
                this.$utils.postApiRequest({url: '/auth/phone', params, isShowError: false}).subscribe(res => {
                    if (res.ok) {
                        this.isRegisterSuccess = true
                    } else {
                        if (res.message.includes('phone already exist')) {
                            this.$q.notify({
                                message: 'Пользователь с таким номером телефона уже зарегистрирован',
                                type: 'negative',
                                position: 'top-right'
                            })
                        } else {
                            this.$q.notify({message: res.message, type: 'negative', position: 'top-right'})
                        }
                    }
                })
            },
            checkCode() {
                if (!this.smsCode) {
                    this.$q.notify({
                        message: 'Необходимо заполнить поле с кодом из sms',
                        type: 'negative',
                        position: 'top-right'
                    })
                    return
                }
                // добавляем флаг, что это регистрация нового пользователя, а не авторизация по логину и паролю
                let params = {token: this.smsCode, phone: this.regForm.login.replace(/\D/g, '')}
                this.$utils.postApiRequest({url: '/auth/check_sms_code', params, isShowError: false}).subscribe(res => {
                    if (res.ok) {
                        this.$currentUser.login({user: res.result})
                    } else {
                        this.$q.notify({
                            message: 'Неверный код подтверждения. Попробуйте еще раз.',
                            type: 'negative',
                            position: 'top-right'
                        })
                    }
                })
            },
        },
        mounted() {
            this.flds.map(v => this.$set(this.regForm, v, null))
        }
    }
</script>
