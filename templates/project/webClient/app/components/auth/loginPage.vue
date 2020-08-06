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
        <email-btn/>
      </div>
    </div>

  </div>
</template>

<script>
    import emailBtn from './email/emailAuthBtn'

    export default {
        components: {emailBtn},
        data() {
            return {}
        },
    }
</script>
