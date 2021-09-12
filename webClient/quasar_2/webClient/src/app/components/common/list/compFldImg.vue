<template>
  <div>
    <div class="q-gutter-md row items-start" style="position: relative">
      <comp-stat-img-src :src="localSrc">
        <div v-if="label" class="absolute-top-left text-subtitle2">
          {{label}}
        </div>
      </comp-stat-img-src>
      <q-btn outline round icon="add" color="grey" @click="isShowDialog = true" class="absolute-top-right all-pointer-events">
        <q-tooltip>Загрузить фото</q-tooltip>
      </q-btn>
      <q-btn outline round size="sm" icon="delete" color="grey" @click="isShowDeleteDialog=true"
             class="absolute-bottom-right all-pointer-events">
        <q-tooltip>{{$t('message.delete')}} фото</q-tooltip>
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
        :accept="(ext && ext.accept) ? ext.accept : ''"
        :max-file-size="(ext && ext.maxFileSize) ? ext.maxFileSize : 10000000"
        @rejected="rejected"
        @uploaded='uploaded'
        @failed='failed'
        :form-fields="formField"
      />
    </q-dialog>
    <!-- диалог подтверждения удаления   -->
    <q-dialog v-model="isShowDeleteDialog" persistent>
      <q-card>
        <q-card-section class="row items-center">
          <q-avatar rounded icon="warning" color="warning" text-color="white"/>
          <span class="q-ml-sm">{{$t('message.delete')}}?</span>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="$t('message.cancel')" v-close-popup/>
          <q-btn flat :label="$t('message.delete')" v-close-popup @click="remove"/>
        </q-card-actions>
      </q-card>
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
      uploadUrl: function () {
        return `${this.$config.apiUrl()}/api/${this.ext.uploadUrl || 'upload_image'}`
      },
      headers: function () {
        const authToken = localStorage.getItem(this.$config.appName)
        return [{name: 'Auth-token', value: authToken}]
      },
      formField: function () {
        let res = [{name: 'tableName', value: this.ext.tableName}, {name: 'tableId', value: this.ext.tableId}]
        if (this.ext.width) res.push({name: 'width', value: this.ext.width})
        if (this.ext.crop) res.push({name: 'crop', value: this.ext.crop})
        return res
      },
    },
    data() {
      return {
        isShowDialog: false,
        isShowDeleteDialog: false,
        localSrc: null,
      }
    },
    watch: {
      fld(v) {
        if (v) this.localSrc = v
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
          // обновляем запись, если указано fldName
          if (this.ext.fldName) {
            this.$utils.postCallPgMethod({
              method: `${this.ext.methodUpdate || this.ext.tableName + '_update'}`,
              params: {id: this.ext.tableId, [this.ext.fldName]: this.localSrc}
            }).subscribe(res => {
            })
          }
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
      rejected(msg) {
        const niceBytes = (x) => {
          const units = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
          let l = 0, n = parseInt(x, 10) || 0
          while (n >= 1024 && ++l) {
            n = n / 1024
          }
          return (n.toFixed(n < 10 && l > 0 ? 1 : 0) + ' ' + units[l])
        }
        let msgText = 'данный файл не соответствует ограничениям'
        if (msg.length > 0 && msg[0].failedPropValidation === 'accept') {
          msgText = `Допустимы только файлы с раширением: ${this.ext.accept} `
        }
        if (msg.length > 0 && msg[0].failedPropValidation === 'max-file-size') {
          let size = niceBytes(this.ext.maxFileSize || 10000000)
          msgText = `Допустимы только файлы не больше: ${size}`
        }
        this.$q.notify({
          color: 'negative',
          position: 'bottom',
          message: msgText,
        })
      },
      remove() {
        // обновляем запись, если указано fldName
        if (this.ext.fldName) {
          this.$utils.postCallPgMethod({
            method: `${this.ext.methodUpdate || this.ext.tableName + '_update'}`,
            params: {id: this.ext.tableId, [this.ext.fldName]: null}
          }).subscribe(res => {
            this.localSrc = null
            this.$emit('update', null)
          })
        } else {
          this.localSrc = null
          this.$emit('update', null)
        }
      },
    },
    mounted() {
      if (!this.ext) {
        throw new Error('compFldFiles missed param: "ext"')
      }
      // if (!this.ext.fldName) {
      //     throw new Error('compFldFiles missed param: "ext.fldName"')
      // }
      if (!this.ext.methodUpdate && !(this.ext.tableId && this.ext.tableName)) {
        throw new Error('compFldFiles missed param: "ext.methodUpdate" OR "ext.tableId" AND "ext.tableName"')
      }
      this.localSrc = this.fld || null
    }
  }
</script>
