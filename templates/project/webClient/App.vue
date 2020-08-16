<template>
  <div id="q-app">
    <!-- страница авторизации  -->
    <auth-comp v-if="!isLoggedIn && !isInLogingProcess"/>

    <!-- страница ожидания подтверждения авторизации   -->
    <waiting-auth-page v-if="isWaitingAuth"/>

    <!-- страница ожидания подтверждения авторизации   -->
    <fired-page v-if="isFired"/>

    <!-- основная страница, после авторизации  -->
    <q-layout view="hHh Lpr lFf" v-if="isLoggedIn && isWorking">

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
            <toolbar-task-btn :counter="taskCounter" @click="toggleTaskList"/>
            <toolbar-message-btn :counter="messageCounter" @click="toggleMsgList"/>
            <current-user-toolbar-menu :currentUser="currentUser" @logout="logout"/>
          </div>

        </q-toolbar>
      </q-header>

      <!-- боковое меню     -->
      <side-menu :leftSide="leftSide" :currentUser="currentUser" @hide="leftSide=false"/>
      <!-- список сообщений     -->
      <message-list ref="messageList" :rightSide="isShowMsgList" :currentUser="currentUser" @hide="isShowMsgList = false" @updateCounter="v => messageCounter = v"/>
        [[if not .Vue.IsHideTaskToolbar]]<!-- список задач     -->
      <task-list ref="taskList" :rightSide="isShowTaskList" :currentUser="currentUser" @hide="isShowTaskList = false" @updateCounter="v => taskCounter = v"/>[[end]]

      <q-page-container>
        <router-view :key="$route.fullPath" :currentUser="currentUser" @reloadMsgList="$refs.messageList ? $refs.messageList.reload() : ''"/>
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
    import toolbarMessageBtn from './app/components/currentUser/messages/toolbarMessageBtn'
    [[if not .Vue.IsHideTaskToolbar]]import toolbarTaskBtn from './app/components/currentUser/tasks/toolbarTaskBtn'[[end]]
    import messageList from './app/components/currentUser/messages/list'
    import taskList from './app/components/currentUser/tasks/list'
    import waitingAuthPage from './app/components/auth/waitingAuthPage'
    import firedPage from './app/components/auth/firedPage'

    export default {
        mixins: [currentUserMixin],
        components: {waitingAuthPage, currentUserToolbarMenu, sideMenu, authComp, toolbarMessageBtn, [[if not .Vue.IsHideTaskToolbar -]]toolbarTaskBtn,[[- end]] messageList, taskList, firedPage},
        data() {
            return {
                leftSide: true,
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
            // подключение к получению обновлений с сервера по sse
            const token = localStorage.getItem(this.$config.appName)
            let source = new EventSource(`${this.$config.apiUrl()}/api/sse?authToken=${token}`)
            source.onmessage = (e) => {
                if (typeof e.data === 'string') {
                    let res = JSON.parse(e.data)
                    // вариант получения message для вывода в боковом правом списке
                    if (res.sse_type === 'message') {
                        if (res.data && typeof res.data === 'string') res.data = JSON.parse(res.data)
                        if (res.options && typeof res.options === 'string') res.options = JSON.parse(res.options)
                        this.$refs.messageList.newMessage(res)
                    }
                    // вариант получения task для вывода в боковом правом списке
                    if (res.sse_type === 'task') {
                        if (res.data && typeof res.data === 'string') res.data = JSON.parse(res.data)
                        this.$userTasks.addTask(res)
                    }
                }
            }
        }
    }
</script>
