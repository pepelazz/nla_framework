<template>
  <q-btn round dense flat color="grey-8" icon="access_alarm" @click="$emit('click')">
    <q-badge v-if="counter>0" color="red" text-color="white" floating>
      {{counter}}
    </q-badge>
    <q-tooltip>Задачи</q-tooltip>
  </q-btn>
</template>

<script>
    export default {
        props: ['counter'],
    }
</script>
