<template>
  <div>
  <q-item >
    <q-avatar square>
      <q-icon :name="item.type" :color="item.type"/>
    </q-avatar>
    <q-item-section>
      <q-item-label>{{item.title}}</q-item-label>
      <q-item-label caption>{{formatDate(item.created_at)}}</q-item-label>
      <q-item-label>
        <q-btn v-if="item.data.message" label="детали" size="xs" flat
               @click="showDetailDialog(item)"/>
      </q-item-label>
    </q-item-section>
    <q-item-section side>
      <div class="text-grey-8">
        <q-btn size="12px" flat dense round icon="done" @click="markAsRead(item.id)"/>
      </div>
    </q-item-section>
  </q-item>
    <!-- диалог с подробностями   -->
    <q-dialog
      v-model="isShowDetailDialog"
    >
      <q-card style="width: 700px; max-width: 80vw;">
        <q-card-section>
          <div class="text-h6">{{detailDialogItem.title}}</div>
        </q-card-section>

        <q-card-section class="q-pt-none" v-if="detailDialogItem.data">
          <span v-html="detailDialogItem.data.message"></span>
        </q-card-section>

        <q-card-actions align="right" class="bg-white text-teal">
          <q-btn flat label="OK" v-close-popup/>
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
  import moment from 'moment'
  export default {
      props: ['item'],
      data() {
          return {
              isShowDetailDialog: false,
              detailDialogItem: {},
          }
      },
      methods: {
          showDetailDialog(item) {
              this.detailDialogItem = item
              this.isShowDetailDialog = true
          },
          formatDate(d) {
              return moment(d).format('DD/MM hh:mm')
          },
          markAsRead(id) {
              this.$utils.postCallPgMethod({method: 'message_mark_as_read', params: {id}}).subscribe(res => {
                  if (res.ok) {
                      this.$emit('markAsRead', id)
                  }
              })
          },
      },
  }
</script>
