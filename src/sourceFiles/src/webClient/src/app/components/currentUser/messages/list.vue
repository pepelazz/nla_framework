<template>
  <div>
    <q-drawer :value="rightSide" side="right" bordered @hide="$emit('hide')">
      <q-list separator>
        <q-item-label header>
          Сообщения
        </q-item-label>
        <div style="position: absolute; top: 10px; right: 10px">
          <q-btn round flat color="secondary" icon="refresh" size="sm" @click="reload"/>
        </div>
        <q-separator/>
        <q-item v-for="item in listForRender" :key="item.id">
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
      </q-list>
      <q-separator/>
    </q-drawer>

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
        props: ['currentUser', 'rightSide'],
        computed: {
            listForRender: function () {
                return this.list.filter(v => !v.is_read)
            }
        },
        data() {
            return {
                list: [],
                isShowDetailDialog: false,
                detailDialogItem: {},
            }
        },
        methods: {
            newMessage(msg) {
                if (this.list.findIndex(v => msg.id === v.id) === -1) {
                    if (!msg.type) msg.type = 'info'
                    this.list.unshift(msg)
                    this.showNotifyMsg()
                    this.updateCounter()
                }
            },
            reload() {
                this.$utils.postCallPgMethod({
                    method: 'message_list',
                    params: {is_read: false, order_by: 'created_at desc'}
                }).subscribe(res => {
                    if (res.ok) {
                        this.list = res.result.map(v => {
                            if (!v.type) v.type = 'info'
                            return v
                        })
                        if (this.list.length > 0) {
                            this.updateCounter()
                            this.showNotifyMsg()
                        }
                    }
                })
            },
            updateCounter() {
                this.$emit('updateCounter', this.list.filter(v => !v.is_read).length)
            },
            markAsRead(id) {
                this.$utils.postCallPgMethod({method: 'message_mark_as_read', params: {id}}).subscribe(res => {
                    if (res.ok) {
                        const i = this.list.findIndex(v => v.id === id)
                        this.list[i].is_read = true
                        this.updateCounter()
                    }
                })
            },
            showDetailDialog(item) {
                this.detailDialogItem = item
                this.isShowDetailDialog = true
            },
            formatDate(d) {
                return moment(d).format('DD/MM hh:mm')
            },
            showNotifyMsg() {
                this.$q.notify({
                    position: 'top-right',
                    message: 'У Вас новые непрочитанные сообщения',
                    avatar: 'https://image.flaticon.com/icons/svg/945/945202.svg',
                    timeout: 1000,
                    // actions: [{icon: 'close', color: 'white'}]
                })
            },
        },
        mounted() {
            if (this.currentUser.id) {
                this.reload()
            }
        }
    }
</script>
