<template>
  <div class="layout">
    <div>
      [[if .Config.Logo]]
      <div class="row justify-center" style="margin-top: 20px; margin-bottom: 10px">
        <img
          src="[[.Config.Logo]]"
          alt="" style="width: auto; max-height: 100px">
      </div>
      [[- end]]
      <!--Кнопки авторизации-->
      <div class="row justify-center q-mt-lg">
        <q-banner class="bg-grey-3">
          <template v-slot:avatar>
            <q-avatar rounded>
              <img src="image/waitingAuth.png">
            </q-avatar>
          </template>
          Ожидаем подтверждения Вашей авторизации
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
