<template>
  <div class="row justify-center">
    <!-- создание чата  -->
    <q-btn v-if="isShowNewChatBtn" round color="secondary" @click="createChat">
      <q-avatar size="52px">
        <img src="https://image.flaticon.com/icons/svg/3081/3081489.svg">
      </q-avatar>
      <q-tooltip>создать чат</q-tooltip>
    </q-btn>

    <q-card v-if="chat" style="width: 100%; max-width: 400px">
      <q-bar>
        <q-icon name="chat" />
        <div>чат</div>

      </q-bar>
      <q-scroll-area
        ref="chatScroll"
        :thumb-style="thumbStyle"
        :bar-style="barStyle"
        style="height: 400px; max-width: 400px;"
      >
        <q-card-section>
          <q-chat-message v-for="msg in msgList" :key="msg.id"
                          :name="msg.fullname"
                          :avatar="msg.avatar"
                          :text="[msg.title]"
                          :stamp="$utils.formatPgDateTime(msg.updated_at)"
                          sent
          />
        </q-card-section>
      </q-scroll-area>
      <q-card-section>
        <q-input label="написать сообщение" outlined v-model="newMsg" autogrow>
          <template v-slot:append>
            <q-icon v-if="newMsg && newMsg !== ''" name="send" @click="send" class="cursor-pointer" />
          </template>
        </q-input>
      </q-card-section>
    </q-card>
  </div>
</template>

<script>
    export default {
        props: ['table_name', 'table_id'],
        data() {
            return {
                chat: null,
                isShowNewChatBtn: false,
                msgList: [],
                newMsg: null,
                thumbStyle: {
                    right: '4px',
                    borderRadius: '5px',
                    backgroundColor: '#027be3',
                    width: '5px',
                    opacity: 0.75
                },
                barStyle: {
                    right: '2px',
                    borderRadius: '9px',
                    backgroundColor: '#027be3',
                    width: '9px',
                    opacity: 0.2
                }
            }
        },
        methods: {
            createChat() {
                this.isShowNewChatBtn = false
                this.$utils.postCallPgMethod({method: 'chat_update', params: {id: -1, table_name: this.table_name, table_id: this.table_id}}).subscribe(res => {
                    if (res.ok) {
                        this.reload()
                        this.newMsg = null
                    }
                })
            },
            send() {
                this.$utils.postCallPgMethod({method: 'chat_message_update', params: {id: -1, chat_id: this.chat.id, title: this.newMsg}}).subscribe(res => {
                    if (res.ok) {
                        this.reload()
                        this.newMsg = null
                    }
                })
            },
            reload() {
                this.$utils.postCallPgMethod({method: 'chat_for_table_id', params: {table_name: this.table_name, table_id: this.table_id}, isShowError: false}).subscribe(res => {
                    if (res.ok) {
                        this.isShowNewChatBtn = false
                        if (res.result.length > 0) {
                            this.chat = res.result[0]
                            this.msgList = res.result[0].message_list || []
                            this.msgList = this.msgList.map(v => {
                                if (v.avatar) {
                                    v.avatar = `${this.$config.apiUrl()}${v.avatar}`
                                    console.log('v.avatar:', v.avatar)
                                } else {
                                    v.avatar = 'https://www.svgrepo.com/show/95333/avatar.svg'
                                }
                                return v
                            })
                            this.$nextTick(() => {
                                const scrollArea = this.$refs.chatScroll
                                const scrollTarget = scrollArea.getScrollTarget()
                                const duration = 100; // ms - use 0 to instant scroll
                                scrollArea.setScrollPosition(scrollTarget.scrollHeight, duration)
                            })
                        }
                    }
                })

            }
        },
        mounted() {
            this.reload()
            setTimeout(() => {
                if (!this.chat) {
                    this.isShowNewChatBtn = true
                }
            }, 400)
        }
    }
</script>
