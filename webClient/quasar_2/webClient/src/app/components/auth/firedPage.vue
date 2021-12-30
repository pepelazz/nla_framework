<template>
  <div class="layout">
    <div>
      <div class="row justify-center" style="margin-top: 20px; margin-bottom: 10px">
        <img
          src="[[appLogoSrc]]"
          alt="" style="width: auto; max-height: 100px">
      </div>
      <!--Кнопки авторизации-->
      <div class="row justify-center q-mt-lg">
        <q-banner class="bg-grey-3">
          <template v-slot:avatar>
            <q-avatar rounded>
              <img src="image/fired.png">
            </q-avatar>
          </template>
          Ваш аккаунт заблокирован. Обратитесь к администратору.
        </q-banner>
      </div>
    </div>

  </div>
</template>

<script>
    import config from "src/app/plugins/config";
    export default {
        data() {
            return {}
        },
        mounted() {
          // стираем токен
          localStorage.removeItem(config.appName)
        }
    }
</script>
