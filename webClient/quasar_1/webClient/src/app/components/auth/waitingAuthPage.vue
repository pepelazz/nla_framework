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
    export default {
        data() {
            return {}
        },
    }
</script>
