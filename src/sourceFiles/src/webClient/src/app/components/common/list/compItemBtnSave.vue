<template>
  <div class="q-py-md">
    <q-btn color="secondary" label="сохранить" class="q-mr-sm" @click="$emit('save')"/>
    <q-btn outline color="secondary" label="отмена" @click="$emit('cancel')"/>
  </div>
</template>
