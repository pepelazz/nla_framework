<template>
  <div>
    <q-bar class="bg-secondary text-white shadow-2">
      <div>{{label}}</div>
      <q-space/>
      <q-btn dense flat icon="expand_less" v-if="isShowList && list.length > 0" @click="isShowList = false"/>
      <q-btn dense flat icon="expand_more" v-if="!isShowList && list.length > 0" @click="isShowList = true"/>
      <q-btn dense flat icon="add" @click="isShowDialog = true" v-if="!readonly"/>
    </q-bar>

    <!-- список   -->
    <q-list bordered separator v-if="isShowList">
      <q-item v-for="item in filteredList" :key="item.filename">
        <q-item-section avatar @click="downloadFile(item)">
          <q-avatar rounded>
            <img src="https://image.flaticon.com/icons/svg/1037/1037308.svg">
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label>{{item.filename}}</q-item-label>
        </q-item-section>
        <q-item-section side v-if="!readonly">
          <q-btn flat round icon="delete" size="sm" @click="showDeleteDialog(item.filename)">
            <q-tooltip>{{$t('message.delete')}}</q-tooltip>
          </q-btn>
        </q-item-section>
      </q-item>
    </q-list>

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
        :form-fields="[{name: 'tableName', value: ext.tableName}, {name: 'tableId', value: ext.tableId},]"
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
    import axios from 'axios'
    export default {
        props: {
            fld: {},
            fldName: {
                default: 'files',
            },
            label: {},
            readonly: null,
            icon: null,
            vif: {
                default: true,
            },
            ext: null,
        },
        computed: {
            uploadUrl() {
                return `${this.$config.apiUrl()}/api/upload_file`
            },
            headers: function () {
                const authToken = localStorage.getItem(this.$config.appName)
                return [{name: 'Auth-token', value: authToken}]
            },
            filteredList: function () {
                return this.list.filter(v => !v.deleted)
            }
        },
        data() {
            return {
                isShowDialog: false,
                isShowList: true,
                list: [],
                isShowDeleteDialog: false,
                selectedForDeleteFilename: null,
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
                    this.$refs.uploader.reset()
                    this.isShowDialog = false
                    // обновляем список если это файл с новым именем или с существующим, но помечен как удаленный
                    if (!this.list.find(v => v.filename === res.result.filename && !v.deleted)) {
                        this.list.push(res.result)
                        this.$utils.postCallPgMethod({
                            method: `${this.ext.tableName}_update`,
                            params: {id: this.ext.tableId, [this.fldName]: this.list}
                        }).subscribe(res => {
                            this.$emit('update', this.list)
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
            downloadFile(item) {
                // через vue.http.post почему то корректно  xsl файл не скачивается
                axios({
                    url: `${this.$config.apiUrl()}${item.url}`,
                    method: 'GET',
                    headers: {'Auth-token': localStorage.getItem(this.$config.appName)},
                    responseType: 'blob', // important
                }).then(response => {
                    const url = window.URL.createObjectURL(new Blob([response.data]));
                    const link = document.createElement('a');
                    link.href = url;
                    link.setAttribute('download', item.filename);
                    document.body.appendChild(link);
                    link.click();
                }).catch((err) => {
                    let msg = err.response.data
                    if (err.response.data) {
                        msg = err.response.data.message
                        if (msg === 'not found') msg = 'файл с такой ссылкой не найден'
                    }
                    this.$q.notify({message: msg, type: 'negative', position: 'top-right'})
                })
            },
            showDeleteDialog(filename) {
                this.selectedForDeleteFilename = filename
                this.isShowDeleteDialog = true
            },
            remove() {
                let i = this.list.findIndex(v => v.filename === this.selectedForDeleteFilename && !v.deleted)
                // помечаем файл как удааленный
                this.list[i].deleted = true
                const fileToken = this.list[i].url.split('/', -1).slice(-1)
                let item = this.list.splice(i, 1)
                this.list.push(item[0])
                // удаляем файл физически
                this.$utils.postApiRequest({url: `/api/remove_file/${fileToken}`}).subscribe(res => {
                    if (res.ok) {
                        this.selectedForDeleteFilename = null
                    }
                })
                // обновляем список файлов в самой записи
                this.$utils.postCallPgMethod({
                    method: `${this.ext.tableName}_update`,
                    params: {id: this.ext.tableId, [this.fldName]: this.list}
                }).subscribe(res => {
                    this.$emit('update', this.list)
                })
            },
        },
        mounted() {
            if (!this.ext) {
                throw new Error('compFldFiles missed param: "ext"')
            }
            if (!this.ext.tableName) {
                throw new Error('compFldFiles missed param: "ext.tableName"')
            }
            if (!this.ext.tableId) {
                throw new Error('compFldFiles missed param: "ext.tableId"')
            }
            this.list = this.fld || []
            // добавляем axios interceptors для обработки ошибки при скачивании файла
            axios.interceptors.response.use(
                response => { return response; },
                error => {
                    if (
                        error.request.responseType === 'blob' &&
                        error.response.data instanceof Blob &&
                        error.response.data.type &&
                        error.response.data.type.toLowerCase().indexOf('json') !== -1
                    )
                    {
                        return new Promise((resolve, reject) => {
                            let reader = new FileReader();
                            reader.onload = () => {
                                error.response.data = JSON.parse(reader.result);
                                resolve(Promise.reject(error));
                            };

                            reader.onerror = () => {
                                reject(error);
                            };

                            reader.readAsText(error.response.data);
                        });
                    };

                    return Promise.reject(error);
                }
            );
        }
    }
</script>
