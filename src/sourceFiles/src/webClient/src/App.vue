<template>
  <div id="q-app">
    <!-- страница авторизации  -->
    <auth-comp v-if="!isLoggedIn && !isInLogingProcess"/>

    <!-- основная страница, после авторизации  -->
    <q-layout view="hHh Lpr lFf" v-if="isLoggedIn">

      <q-header reveal elevated class="bg-white text-grey-8 q-py-xs">
        <q-toolbar>
          <q-btn dense flat round icon="menu" @click="leftSide = !leftSide"/>

          <!-- лого -->
          <q-btn flat no-caps no-wrap class="q-ml-xs" @click="$router.push('/')">
<!--            <q-avatar size="26px">-->
<!--              <img src="https://www.defly.ru/website/defly/template/images/logo.png">-->
<!--            </q-avatar>-->
            <q-toolbar-title shrink class="text-weight-bold">
              {{$config.uiAppName}}
            </q-toolbar-title>
          </q-btn>

          <q-space/>

          <!-- аватарка и меню пользователя -->
          <div class="q-gutter-sm row items-center no-wrap">
            <current-user-toolbar-menu :currentUser="currentUser" @logout="logout"/>
          </div>

        </q-toolbar>
      </q-header>

      <!-- боковое меню     -->
      <side-menu :leftSide="leftSide" :currentUser="currentUser"/>

      <q-page-container>
        <router-view :key="$route.fullPath"/>
      </q-page-container>

    </q-layout>
  </div>
</template>

<script>
    import currentUserMixin from './app/mixins/currentUser'
    // import authPage from './app/components/auth'
    import currentUserToolbarMenu from './app/components/currentUser/toolbarMenu'
    import sideMenu from './app/components/sidemenu/index.vue'
    import authComp from './app/components/auth/index'

    export default {
        mixins: [currentUserMixin],
        components: {currentUserToolbarMenu, sideMenu, authComp},
        data() {
            return {
                leftSide: true,
            }
        },
        mounted() {
            this.$currentUser.login()
        }
    }
</script>

<style>
</style>
