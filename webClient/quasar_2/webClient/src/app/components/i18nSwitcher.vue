<template>
  <q-select
    v-model="locale"
    :options="localeOptions"
    dense
    borderless
    emit-value
    map-options
    options-dense
  />
</template>

<script>
  import { ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'

  export default {
    setup () {
      const { locale } = useI18n({ useScope: 'global' })
      if (localStorage.getItem("i18n")) locale.value = localStorage.getItem("i18n")

      watch(locale, ()=> localStorage.setItem('i18n', locale.value))

      return {
        locale,
        localeOptions: [
          { value: 'ru-RU', label: 'RU' },
          { value: 'en-US', label: 'EN' },
        ]
      }
    }
  }
</script>
