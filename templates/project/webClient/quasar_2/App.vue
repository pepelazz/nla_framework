<template>
  <div>
    <!-- страница авторизации  -->
    <auth-comp v-if="!isLoggedIn && !isInLogingProcess"/>
    <!-- страница ожидания подтверждения авторизации   -->
    <waiting-auth-page v-if="isWaitingAuth"/>

    <!-- страница ожидания подтверждения авторизации   -->
    <fired-page v-if="isFired"/>

    <!-- основная страница, после авторизации  -->
    <q-layout view="hHh lpR fFf" v-if="isLoggedIn && isWorking">

      <q-header elevated class="bg-white text-grey-8 q-py-xs">
        <q-toolbar>
          <q-btn dense flat round icon="menu" @click="leftSide = !leftSide"/>

          <!-- лого -->
          <q-btn flat no-caps no-wrap class="q-ml-xs" @click="$router.push('/')">
            <!--            <q-avatar size="26px">-->
            <!--              <img src="https://www.defly.ru/website/defly/template/images/logo.png">-->
            <!--            </q-avatar>-->


            [[if not .Vue.UiAppLogoOnly]]
            <q-toolbar-title shrink class="text-weight-bold">
              {{$config.uiAppName}}
            </q-toolbar-title>
            [[- end]]
            [[if .Vue.UiAppLogoOnly]]
            <q-toolbar-title shrink class="text-weight-bold">
              [[.Vue.UiAppLogoOnly]]
            </q-toolbar-title>
            [[- end]]
          </q-btn>

          <q-space/>

          [[- if .I18n.IsExist]]
          <i18n-switcher/>
          [[- end]]

          <!-- аватарка и меню пользователя -->
          <div class="q-gutter-sm row items-center no-wrap">


            <current-user-toolbar-menu :currentUser="currentUser" @logout="logout"/>
          </div>

        </q-toolbar>
      </q-header>

      <!-- боковое меню     -->
      <side-menu :leftSide="leftSide" :currentUser="currentUser" @hide="leftSide=false"/>

      <q-page-container>
        <router-view :key="$route.fullPath" :currentUser="currentUser"/>
      </q-page-container>

    </q-layout>
  </div>
</template>

<script>
  import currentUserMixin from './app/mixins/currentUser'
  import currentUserToolbarMenu from './app/components/currentUser/toolbarMenu'
  import sideMenu from './app/components/sidemenu/index.vue'
  import authComp from './app/components/auth/index'

  import waitingAuthPage from './app/components/auth/waitingAuthPage'
  import firedPage from './app/components/auth/firedPage'
  [[- if .I18n.IsExist]]
  import i18nSwitcher from "src/app/components/i18nSwitcher"
  [[- end]]

  export default {
    mixins: [currentUserMixin],
    components: {authComp, currentUserToolbarMenu, sideMenu, waitingAuthPage, firedPage, [[- if .I18n.IsExist]]i18nSwitcher,[[- end]]},
    data() {
      return {
        leftSide: false,
        isShowMsgList: false,
        isShowTaskList: false,
        messageCounter: 0,
        taskCounter: 0,
      }
    },
    methods: {
      toggleTaskList() {
        if (this.isShowTaskList) {
          this.isShowTaskList = false
        } else {
          if (this.isShowMsgList) {
            this.isShowMsgList = false
            this.$nextTick(() => this.isShowTaskList = true)
          } else {
            this.isShowTaskList = true
          }
        }
      },
      toggleMsgList() {
        if (this.isShowMsgList) {
          this.isShowMsgList = false
        } else {
          if (this.isShowTaskList) {
            this.isShowTaskList = false
            this.$nextTick(() => this.isShowMsgList = true)
          } else {
            this.isShowMsgList = true
          }
        }
      },
      hideRightSidebar() {
        this.isShowMsgList = false
        this.isShowTaskList = false
      }
    },
    mounted() {
      this.$currentUser.login()
    }
  }
</script>
