<template>
  <q-input outlined :label="label" v-model="localItem" autogrow/>
</template>

<script>
    export default {
        props: {
            label: {},
            fld: {},
            ext: null, // дополнительные параметры
            readonly: null,
            columnClass: {
                default: 'col-xs-12 col-sm-6 col-md-4'
            },
            icon: null,
            vif: {
                default: true,
            },
        },
        data() {
            return {
                localItem: null,
            }
        },
        watch: {
            localItem: function () {
                this.update()
            },
        },
        methods: {
            update() {
                this.$emit('update', {text: this.localItem})
            },
        },
        mounted() {
            if (this.fld) this.localItem = this.fld.text
        }
    }
</script>
