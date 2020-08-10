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
        <component v-for="item in listForRender" :key="item.id" :is="item.template" :item="item" @markAsRead="markAsRead"></component>
      </q-list>
      <q-separator/>
    </q-drawer>
  </div>
</template>

<script>
    import moment from 'moment'
    import defaultTmpl from './msgTemplate/default.vue'
    [[range .Vue.MessageTmpls]]
    import [[.CompName]] from '[[.CompPath]]'
    [[- end]]

    export default {
        props: ['currentUser', 'rightSide'],
      components: {defaultTmpl [[range .Vue.MessageTmpls]], [[.CompName]] [[- end]]},
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
                    msg.template = msg.options && msg.options.template ? msg.options.template : 'defaultTmpl'
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
                            v.template = v.options && v.options.template ? v.options.template : 'defaultTmpl'
                            return v
                        })
                        this.updateCounter()
                        if (this.list.length > 0) {
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
