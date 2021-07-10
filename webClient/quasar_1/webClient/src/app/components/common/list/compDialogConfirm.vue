<template>
  <q-dialog v-model="isShowDialog" persistent>
    <q-card>
      <q-card-section class="row items-center">
        <q-avatar v-if="icon" :icon="icon" color="primary" text-color="white" />
        <span class="q-ml-sm">{{label}}</span>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Отмена" color="primary" v-close-popup />
        <q-btn flat label="Ок" color="primary" v-close-popup @click="$emit('update', item)"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script>
  export default {
    props: ['label', 'icon'],
    data() {
      return {
        isShowDialog: false,
        item: null,
      }
    },
    methods: {
      open(item) {
        this.isShowDialog = true
        this.item = item
      }
    }
  }
</script>
