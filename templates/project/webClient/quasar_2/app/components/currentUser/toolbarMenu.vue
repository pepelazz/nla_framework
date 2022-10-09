<template>
  <q-btn round flat>
    <q-avatar rounded size="26px">
      <comp-stat-img-src v-if="currentUser.avatar" :src="currentUser.avatar"/>
      <img v-else src="https://www.svgrepo.com/show/95333/avatar.svg">
    </q-avatar>
    <q-menu
      transition-show="flip-right"
      transition-hide="flip-left"
      auto-close>
      <q-list dense>
        <q-item class="GL__menu-link-signed-in">
          <q-item-section>
            <div><strong>{{currentUser.fullname}}</strong></div>
          </q-item-section>
        </q-item>
        <q-separator/>
        <q-item clickable class="GL__menu-link">
          <q-item-section @click="$router.push('/profile')">{{$t('message.edit')}}</q-item-section>
        </q-item>
        <q-separator/>
        <q-item clickable class="GL__menu-link">
          <q-item-section @click="$emit('logout')">{{$t('profile.exit')}}</q-item-section>
        </q-item>
        [[if .Vue.Theme.IsDarkThemeExist ]]
        <q-separator />
        <q-item clickable class="GL__menu-link">
          <dark-mode-toggle-comp />
        </q-item>
        [[- end]]
      </q-list>
    </q-menu>
  </q-btn>
</template>

<script>
    [[if .Vue.Theme.IsDarkThemeExist ]]
    import {setCssVar} from "quasar";
    import darkModeToggleComp from "./comp/darkModeToggleComp";
    [[- end]]
    export default {
        props: ['currentUser'],
        [[if .Vue.Theme.IsDarkThemeExist ]]
        components: {darkModeToggleComp},
        [[- end]]
        data() {
            return {}
        },
        [[if .Vue.Theme.IsDarkThemeExist ]]
        mounted() {
          const isDarkMode = localStorage.getItem('isDarkMode') == 'true' || false
          this.$q.dark.set(isDarkMode)
          if (isDarkMode) {
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
          if (!isDarkMode) {
            setCssVar('primary', '#1976D2')
            setCssVar('secondary', '#26A69A')
            setCssVar('accent', '#9C27B0')
            setCssVar('positive', '#21BA45')
            setCssVar('negative', '#C10015')
            setCssVar('info', '#31CCEC')
            setCssVar('warning', '#F2C037')
          }
        }
        [[- end]]
    }
</script>

[[if .Vue.Theme.IsDarkThemeExist ]]
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
[[- end]]
