<template>
  <div>
    <login-page v-if="isRedirectToLoginPage"/>
    <div v-else>
      <router-view/>
    </div>
  </div>
</template>

<script>
    import loginPage from './loginPage'
    export default {
        components: {loginPage},
        computed: {
            isRedirectToLoginPage: function () {
                return !['/check_user_email', '/email_auth_recover_password'].includes(this.$route.path)
            }
        }
    }
</script>
