<template>
  <div class="layout">
    <div>
      <div class="row justify-center" style="margin-top: 20px; margin-bottom: 10px">
        <img
          src="https://c.radikal.ru/c19/2001/d3/4563ca8eda40.jpg"
          alt="" style="width: auto; max-height: 100px">
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
