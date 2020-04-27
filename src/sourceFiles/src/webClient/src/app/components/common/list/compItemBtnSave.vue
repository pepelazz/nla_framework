<template>
  <div class="q-py-md">
    <q-btn color="secondary" :label="$t('message.save')" class="q-mr-sm" @click="$emit('save')"/>
    <q-btn outline color="secondary" :label="$t('message.cancel')" @click="$emit('cancel')"/>
  </div>
</template>
