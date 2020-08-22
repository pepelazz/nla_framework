<template>
  <div class="layout">
    <div>
      <div class="row justify-center" style="margin-top: 20px; margin-bottom: 10px">
        [[if .Config.Logo]]<img src="[[.Config.Logo]]" alt="" style="width: auto; max-height: 100px">
        [[else]]
        <h5>[[.Vue.UiAppName]]</h5>
        [[- end]]
      </div>
      <!--Кнопки авторизации-->
      <div class="row justify-center q-mt-lg">
        [[if .Config.Auth.ByEmail]]<email-btn/>[[end]]
        [[if .Config.Auth.ByPhone]]<phone-btn/>[[end]]

      </div>
    </div>

  </div>
</template>

<script>
  [[if .Config.Auth.ByEmail]]import emailBtn from './email/emailAuthBtn'[[end]]
  [[if .Config.Auth.ByPhone]]import phoneBtn from './phone/phoneAuthBtn'[[end]]

    export default {
      [[if .Config.Auth.ByEmail]]components: {emailBtn},[[end]]
      [[if .Config.Auth.ByPhone]]components: {phoneBtn},[[end]]
        data() {
            return {}
        },
    }
</script>
