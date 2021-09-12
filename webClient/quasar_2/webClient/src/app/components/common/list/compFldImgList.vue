<template>
  <div>
    <p class="text-caption">{{label}}</p>
    <div class="q-gutter-md row items-start">
      <comp-stat-img-src v-for="(item, index) in list" :key="item.file"
                         :src="item.file" @error="v=> imgSrcError(item.file, v)"
                         style="width: 150px"
      >
        <a :href="imgUrl(item.file)" target="_blank">
          <q-btn flat round size="sm" icon="link" color="white" class="absolute-top-right all-pointer-events"/>
        </a>
        <q-btn v-if="!readonly" flat round size="sm" icon="delete" color="white" @click="showDeleteDialog(item.file)" class="absolute-bottom-right all-pointer-events">
          <q-tooltip>{{$t('message.delete')}} фото</q-tooltip>
        </q-btn>
        <q-btn v-if="!readonly && index>0" flat round size="sm" icon="keyboard_backspace" color="white" @click="moveLeft(item.file, index)" class="absolute-bottom-left all-pointer-events">
          <q-tooltip>Поменять местами</q-tooltip>
        </q-btn>
      </comp-stat-img-src>
      <q-btn v-if="!readonly" flat round icon="add" size="sm" @click="isShowDialog = true">
        <q-tooltip>Добавить фото</q-tooltip>
      </q-btn>
      <q-btn v-if="!readonly && ext && ext.canAddUrls" size="sm" flat round icon="add" @click="isShowDialogAddUrl = true">
        <q-tooltip>Добавить ссылку</q-tooltip>
      </q-btn>
    </div>

    <!-- диалог добавления   -->
    <q-dialog v-model="isShowDialog">
      <q-uploader
        ref="uploader"
        label="Выберите файл для загрузки"
        multiple
        :url="uploadUrl"
        :headers='headers'
        :accept="(ext && ext.accept) ? ext.accept : ''"
        :max-file-size="(ext && ext.maxFileSize) ? ext.maxFileSize : 10000000"
        @rejected="rejected"
        @uploaded='uploaded'
        @failed='failed'
        @finish ='finish'
        :form-fields="formField"
      />
    </q-dialog>

    <!-- диалог добавления url ссылку на фото  -->
    <q-dialog v-model="isShowDialogAddUrl">
      <q-card style="width: 300px" class="q-px-sm q-pb-md">
        <q-img :src="newImgUrl ? newImgUrl : 'https://www.cowgirlcontractcleaning.com/wp-content/uploads/sites/360/2018/05/placeholder-img-4.jpg'" />
        <q-card-section>
          <q-input v-model="newImgUrl" label="ссылка на фото"/>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('message.cancel')" v-close-popup/>
          <q-btn flat label="Ок" v-close-popup @click="addImgUrl"/>
        </q-card-actions>
      </q-card>
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
      uploadUrl: function() {
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
      imgUrl: function () {
        return (src) => src.includes('http') ? src : `${this.$config.apiUrl()}${src}`
      },
    },
    data() {
      return {
        isShowDialog: false,
        list: [],
        isShowDeleteDialog: false,
        isShowDialogAddUrl: false,
        selectedForDeleteFilename: null,
        newImgUrl: null,
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
          this.list.push(res.result)
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
      showDeleteDialog(filename) {
        this.selectedForDeleteFilename = filename
        this.isShowDeleteDialog = true
      },
      imgSrcError(filename, msg) {
        this.$q.notify({
          color: 'negative',
          position: 'bottom',
          message: `ошибка: url "${filename}" не является ссылкой на фото `,
        })
        this.selectedForDeleteFilename = filename
        this.remove()
      },
      moveLeft(file, i) {
        let origin = this.list[i]
        this.list[i] = this.list[i - 1]
        this.$set(this.list, i - 1, origin)
        this.$emit('update', this.list)
      },
      remove() {
        let i = this.list.findIndex(v => v.file === this.selectedForDeleteFilename && !v.deleted)
        // помечаем файл как удааленный
        this.list[i].deleted = true
        let item = this.list.splice(i, 1)
        this.list.push(item[0])
        // обновляем список файлов в самой записи
        this.$utils.postCallPgMethod({
          method: `${this.ext.methodUpdate || this.ext.tableName + '_update'}`,
          params: {id: this.ext.tableId, [this.ext.fldName]: this.list}
        }).subscribe(res => {
          this.$emit('update', this.list)
        })
      },
      // загрузка url ссылки на фото
      addImgUrl: function () {
        this.isShowDialog = false
        this.list.push({file: this.newImgUrl})
        // обновляем запись
        this.$utils.postCallPgMethod({
          method: `${this.ext.methodUpdate || this.ext.tableName + '_update'}`,
          params: {id: this.ext.tableId, [this.ext.fldName]: this.list}
        }).subscribe(res => {
          if (res.ok) {
            this.newImgUrl = null
          }
        })
      },
      finish() {
        this.$refs.uploader.reset()
        this.isShowDialog = false
        this.$emit('update', this.list)
        // обновляем запись
        this.$utils.postCallPgMethod({
          method: `${this.ext.methodUpdate || this.ext.tableName + '_update'}`,
          params: {id: this.ext.tableId, [this.ext.fldName]: this.list}
        }).subscribe(res => {
        })
      },
    },
    mounted() {
      // проверки только если не readonly
      if (!this.readonly) {
        if (!this.ext) {
          throw new Error('compFldFiles missed param: "ext"')
        }
        if (!this.ext.fldName) {
          throw new Error('compFldFiles missed param: "ext.fldName"')
        }
        if (!this.ext.methodUpdate && !(this.ext.tableId && this.ext.tableName)) {
          throw new Error('compFldFiles missed param: "ext.methodUpdate" OR "ext.tableId" AND "ext.tableName"')
        }
      }
      if (this.fld) {
        this.list = this.fld.filter(v => !v.deleted)
      } else {
        this.list = []
      }
    }
  }
</script>
