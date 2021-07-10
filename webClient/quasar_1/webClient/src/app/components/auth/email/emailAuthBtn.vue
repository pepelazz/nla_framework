<template>
  <div>
    <q-btn size='md' outline color='primary' @click='login' :disabled='disabled'>
      <q-icon left size="3em" name="perm_identity"/>
      <div>войти</div>
    </q-btn>
    <q-dialog v-model="modalIsOpened" persistent>
      <q-card style="min-width: 35%">
        <q-toolbar>
          <q-toolbar-title><span class="text-weight-bold">{{dialogTitle}}</span></q-toolbar-title>
          <q-btn flat round dense icon="close" v-close-popup/>
        </q-toolbar>

        <q-card-section>
          <!-- форма логин + пароль -->
          <comp-login-form v-if='!(emailRegister.isNewRegister || emailRegister.isPasswordRecover)'
                           @register="showNewEmailUserRegister"
                           @passwordRecover="showPasswordRecover"/>
          <!--форма регистрации -->
          <comp-register-form v-if="emailRegister.isNewRegister" @cancel="resetRegisterForm"/>
          <!--форма восстановления пароля -->
          <comp-recover-password-form v-if='emailRegister.isPasswordRecover' @cancel="resetPasswordRecover"/>
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
</template>

<style lang="stylus" scoped>
  .bg-facebook
    background: #4267b2
</style>

<script>
  import compLoginForm from './components/compLoginForm'
  import compRegisterForm from './components/compRegisterForm'
  import compRecoverPasswordForm from './components/compRecoverPasswordForm'

  export default {
    props: ['disabled'],
    components: {compLoginForm, compRegisterForm, compRecoverPasswordForm},
    computed: {
      dialogTitle() {
        if (this.emailRegister.isNewRegister) return 'Регистрация'
        if (this.emailRegister.isPasswordRecover) return 'Восстановление пароля'
        return 'Авторизация'
      },
    },
    data() {
      return {
        modalIsOpened: false,
        emailRegister: {
          isNewRegister: false,
          isNewRegisterSuccess: false, // флаг для изменения состояния, когда форма регистрации успешно отправлена
          isPasswordRecover: false,
          isPasswordRecoverSuccess: false,
        },
      }
    },
    methods: {
      login() {
        this.modalIsOpened = true
      },
      showNewEmailUserRegister() {
        this.emailRegister.isNewRegister = true
        this.emailRegister.isNewRegisterSuccess = false
      },
        resetRegisterForm() {
        this.emailRegister.isNewRegister = false
      },
      showPasswordRecover() {
        this.emailRegister.isPasswordRecover = true
      },
      resetPasswordRecover() {
        this.emailRegister.isPasswordRecover = false
      },
    },
  }
</script>
