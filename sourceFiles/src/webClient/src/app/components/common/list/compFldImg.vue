<template>
    <div>
      <div class="q-gutter-md row items-start" style="position: relative">
        <comp-stat-img-src :src="localSrc">
          <div v-if="label" class="absolute-top-left text-subtitle2">
            {{label}}
          </div>
        </comp-stat-img-src>
        <q-btn flat round icon="add" color="white" @click="isShowDialog = true" class="absolute-top-right">
          <q-tooltip>Загрузить фото</q-tooltip>
        </q-btn>
      </div>
      <!-- диалог добавления   -->
      <q-dialog v-model="isShowDialog">
        <q-uploader
          ref="uploader"
          label="Выберите файл для загрузки"
          auto-upload
          :url="uploadUrl"
          :headers='headers'
          @uploaded='uploaded'
          @failed='failed'
          :form-fields="[{name: 'tableName', value: ext.tableName}, {name: 'tableId', value: ext.tableId},]"
        />
      </q-dialog>
    </div>
</template>

<script>
    export default {
        props: {
            fld: {},
            label: {},
            readonly: null,
            icon: null,
            vif: {
                default: true,
            },
            ext: {},
        },
        computed: {
            uploadUrl: function() {
                return `${this.$config.apiUrl()}/api/${this.ext.uploadUrl || 'upload_image'}`
            },
            headers: function () {
                const authToken = localStorage.getItem(this.$config.appName)
                return [{name: 'Auth-token', value: authToken}]
            },
        },
        data() {
            return {
                isShowDialog: false,
                localSrc: null,
            }
        },
        methods: {
            uploaded({xhr: {response}}) {
                const res = JSON.parse(response)
                if (!res.ok) {
                    this.$q.notify({
                        color: 'negative',
                        position: 'bottom',
                        message: res.message,
                    })
                } else {
                    this.localSrc = res.result.file
                    this.$refs.uploader.reset()
                    this.isShowDialog = false
                    this.$emit('update', this.localSrc)
                    // обновляем запись
                    this.$utils.postCallPgMethod({
                        method: `${this.ext.methodUpdate || this.ext.tableName + '_update'}`,
                        params: {id: this.ext.tableId, [this.ext.fldName]: this.localSrc}
                    }).subscribe(res => {
                    })
                }
            },
            failed(msg) {
                let msgText = 'ошибка загрузки'
                if (msg.xhr && msg.xhr.responseText) {
                    let res = JSON.parse(msg.xhr.responseText)
                    if (res.message) msgText = res.message
                }
                this.$q.notify({
                    color: 'negative',
                    position: 'bottom',
                    message: msgText,
                })
            },
        },
        mounted() {
            if (!this.ext) {
                throw new Error('compFldFiles missed param: "ext"')
            }
            if (!this.ext.fldName) {
                throw new Error('compFldFiles missed param: "ext.fldName"')
            }
            if (!this.ext.methodUpdate && !(this.ext.tableId && this.ext.tableName)) {
                throw new Error('compFldFiles missed param: "ext.methodUpdate" OR "ext.tableId" AND "ext.tableName"')
            }
            this.localSrc = this.fld || null
        }
    }
</script>
