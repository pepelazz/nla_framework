<template>
  <div>
    <q-toggle
      v-model="darkMode"
      checked-icon="brightness_2"
      unchecked-icon="wb_sunny"/>
  </div>
</template>

<script>
import {setCssVar, useQuasar} from "quasar";
import {ref, watch, onMounted} from 'vue';

export default {
  name: "darkModeToggleComp",
  setup() {
    const $q = useQuasar()
    const darkMode = ref(localStorage.getItem('isDarkMode') == 'true' || false)
    $q.dark.set(darkMode.value)

    watch(darkMode, val => {
      $q.dark.set(val)
      localStorage.setItem('isDarkMode', val)
      if (val) {
        setCssVar('primary', '#98cfda')
        setCssVar('secondary', '#00675b')
        setCssVar('accent', '#ae99d3')
        setCssVar('positive', '#4db063')
        setCssVar('negative', '#c44823')
        setCssVar('info', '#009688')
        setCssVar('warning', '#F2C037')
        setCssVar('dark', '#494949')
        // setCssVar('dark', '#121212')
        setCssVar('dark-page', '#363636')
        // setCssVar('dark-page', '#000000')
        setCssVar('dark-layer', 'rgba(255, 255, 255, 0.15)')
      }
      if (!val) {
        setCssVar('primary', '#1976D2')
        setCssVar('secondary', '#26A69A')
        setCssVar('accent', '#9C27B0')
        setCssVar('positive', '#21BA45')
        setCssVar('negative', '#C10015')
        setCssVar('info', '#31CCEC')
        setCssVar('warning', '#F2C037')
      }
    })

    return {
      darkMode,
    }
  }
}
</script>

<style lang="scss">
body.body--light {
  .main-header {
    background-color: #fff;
    color: $grey-8;
  }
}

body.body--dark {
  .main-header {
    background-color: #363636;
    color: $grey-6;
  }
}

</style>
