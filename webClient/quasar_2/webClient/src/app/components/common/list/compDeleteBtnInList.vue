<template>
  <div>
    <q-item-section side>
      <q-btn v-if="!item.deleted" flat round icon="delete_forever" @click="isShowDeleteDialog=true">
        <q-tooltip>{{$t('message.delete')}}</q-tooltip>
      </q-btn>
      <q-btn v-if="item.deleted" flat round icon="check_circle_outline" @click="update">
        <q-tooltip>Восстановить</q-tooltip>
      </q-btn>
    </q-item-section>
    <!-- диалог подтверждения удаления   -->
    <q-dialog v-model="isShowDeleteDialog" persistent>
      <q-card>
        <q-card-section class="row items-center">
          <q-avatar rounded icon="warning" color="warning" text-color="white"/>
          <span class="q-ml-sm">{{$t('message.delete')}}?</span>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="$t('message.cancel')" v-close-popup/>
          <q-btn flat :label="$t('message.delete')" v-close-popup @click="update"/>
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
    export default {
        props: ['item', 'updateMethod'],
        data() {
            return {
                list: [],
                isShowDeleteDialog: false,
            }
        },
        methods: {
            update() {
                if (this.updateMethod) {
                    this.$utils.postCallPgMethod({
                        method: this.updateMethod,
                        params: {
                            id: this.item.id,
                            deleted: !this.item.deleted,
                        }
                    }).subscribe(res => {
                        if (res.ok) {
                            // загружаем весь список из базы по новой
                            this.$emit('success')
                        }
                    })
                } else {
                    this.$emit('success')
                }
            }
        },
    }
</script>
