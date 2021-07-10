<!-- компонента для отображения ссылок на фото, загруженных на сервер. В зависимости от режима разработки подставляет адрес сервера-->
<template>
  <q-img :src="imgUrl" @error="v=> $emit('error', v)">
    <slot></slot>
  </q-img>
</template>

<script>
    export default {
        props: ['src'],
        computed: {
            imgUrl: function () {
                if (this.src) {
                    return this.src.includes('http') || this.src.startsWith('image') ? this.src : `${this.$config.apiUrl()}${this.src}`
                } else {
                    return 'https://www.foot.com/wp-content/uploads/2017/03/placeholder.gif'
                }
            }
        },
    }
</script>
