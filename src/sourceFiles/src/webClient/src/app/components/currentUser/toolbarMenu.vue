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
          <q-item-section @click="$router.push('/profile')">редактировать</q-item-section>
        </q-item>
        <q-separator/>
        <q-item clickable class="GL__menu-link">
          <q-item-section @click="$emit('logout')">выйти</q-item-section>
        </q-item>
      </q-list>
    </q-menu>
  </q-btn>
</template>

<script>
    export default {
        props: ['currentUser'],
        data() {
            return {}
        },
    }
</script>
