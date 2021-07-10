<template>
  <div>
    <q-dialog v-model="isShowDialog">
      <q-card style="width: 700px; max-width: 80vw;">
<!--        <q-card-section class="row items-center">-->
<!--          <span class="q-ml-sm">Поиск адреса</span>-->
<!--        </q-card-section>-->

        <q-card-section>
          <comp-dadata-address label="Поиск адреса" @update="update" class="q-mb-md"/>
        </q-card-section>

<!--        <q-card-actions align="right">-->
<!--          <q-btn flat label="Отмена" v-close-popup/>-->
<!--        </q-card-actions>-->
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
    export default {
        props: {
            label: {
                type: String
            },
        },
        data() {
            return {
                isShowDialog: false
            }
        },
        methods: {
            open() {
                this.isShowDialog = true
            },
            update(v) {
                this.$emit('update', v)
                this.isShowDialog = false
            }
        },
    }
</script>
