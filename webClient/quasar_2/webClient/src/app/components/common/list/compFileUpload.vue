<template>
  <div>
    <q-btn v-if="!is_btn_with_label" :icon="iconName" flat round :size="btnSize" @click="isOpenUpload=true">
      <q-tooltip v-if="tooltip">{{tooltip}}</q-tooltip>
    </q-btn>
    <q-btn v-if="is_btn_with_label" :icon="iconName" outline :label="label" :color="btnColor" :size="btnSize" @click="isOpenUpload=true">
      <q-tooltip v-if="tooltip">{{tooltip}}</q-tooltip>
    </q-btn>
    <q-dialog v-model="isOpenUpload">
      <q-card>
        <q-card-section>
          <q-uploader
            ref="uploader"
            :label="labelName"
            auto-upload
            :filter="checkFileType"
            :url="uploadUrl"
            :headers='headers'
            @uploaded='uploaded'
            @failed='failed'
          />
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
  export default {
    props: ['url', 'fileExt', 'icon', 'tooltip', 'label', 'is_btn_with_label', 'size', 'color', 'notifyResultFunc'],
    computed: {
      uploadUrl() {
        return `${this.$config.apiUrl()}/api/${this.url}`
      },
      headers: function () {
        const authToken = localStorage.getItem(this.$config.appName)
        return [{name: 'Auth-token', value: authToken}]
      },
      btnSize() {
        return this.size || 'sm'
      },
      btnColor() {
        return this.color || 'primary'
      },
    },
    data() {
      return {
        isOpenUpload: false,
        iconName: 'cloud_upload',
        labelName: 'Выберите файл для загрузки',
      }
    },
    methods: {
      checkFileType (files) {
        if (this.fileExt && this.fileExt.length > 0) {
          let isRightExt = files.filter(file => this.fileExt.filter(ext => file.name.includes(`.${ext}`)).length > 0)
          if (isRightExt.length === 0) {
            this.$q.notify({type: 'negative', message: `файл должен иметь расширение ${this.fileExt}`})
          }
          return isRightExt
        }
        return files
      },
      uploaded({xhr: {response}}) {
        const res = JSON.parse(response)
        if (!res.ok) {
          this.$q.notify({
            color: 'negative',
            position: 'bottom',
            message: res.message,
          })
        } else {
          this.$refs.uploader.reset()
          this.isOpenUpload = false
          this.$emit('result', res.result)
          this.$emit('reloadList')
          this.$q.notify({
            color: 'positive',
            position: 'bottom',
            message: this.notifyResultFunc ? this.notifyResultFunc(res.result) : res.result,
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
      if (this.icon) this.iconName = this.icon
      if (this.label) this.labelName = this.label
    }
  }
</script>
