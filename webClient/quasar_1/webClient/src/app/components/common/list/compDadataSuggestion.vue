<template>
  <div>
    <q-input
      :label='label'
      v-model="search"
      outlined
      debounce="500"
    >
      <template v-slot:append>
        <q-icon name="search" />
      </template>
    </q-input>
    <q-list bordered separator>
      <q-item v-for="(item, index) in resultList" :key="item.data.hid">
        <q-item-section side>
          <q-item-label>{{index + 1}}
            <q-btn flat round size="sm" icon="add" @click="$emit('add', item)"/>
            <q-btn flat round size="sm" icon="help_outline" @click="openDialogWithDetails(item)"/>
          </q-item-label>
        </q-item-section>
        <q-item-section>
          <q-item-label>{{item.value}}</q-item-label>
          <q-item-label caption>{{item.data.inn}} {{item.data.address.value}}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
    <p class="text-caption" v-if="resultList.length == 0 ? 'ничего не найдено' : null">не найдено</p>

    <!-- диалог с деталями   -->
    <q-dialog v-model="isShowDialog" v-if="selectedItem">
      <q-card style="width: 700px; max-width: 80vw;">
        <q-card-section>
          <div class="text-h6">{{selectedItem.value}}</div>
        </q-card-section>

        <q-separator />

        <q-card-section style="max-height: 50vh" class="scroll">
          <div class="row q-col-gutter-md q-mb-sm" v-for="fldRow in flds">
            <q-input outlined  v-for="fld in fldRow" :key='fld.name'
                      :label="fld.label"
                      :value="fld.value(selectedItem.data)"
                      readonly
                      :class="fld.class"
            />
          </div>
        </q-card-section>

        <q-separator />

        <q-card-actions align="right">
          <q-btn flat label="ok" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
    import config from '../../../plugins/config'
    import {ajax} from 'rxjs/ajax'
    import {catchError, map, take} from 'rxjs/operators'
    import {Notify} from 'quasar'
    import {of} from 'rxjs'

    export default {
        props: {
            label: {
                type: String
            },
        },
        data() {
            return {
                search: '7707083893',
                resultList: [],
                selectedItem: null,
                isShowDialog: false,
                flds: [
                    [
                        {value: (v) => v.inn, label: 'ИНН', class: 'col-xs-12 col-sm-6'},
                        {value: (v) => v.kpp, label: 'КПП', class: 'col-xs-12 col-sm-6'},
                    ],
                    [
                        {value: (v) => v.opf.full, label: 'ОПФ', class: 'col-xs-12 col-sm-6'},
                        {value: (v) => this.$t(`onecContragent.status_${v.state.status}`), label: 'Статус', class: 'col-xs-12 col-sm-6'},
                    ],
                  ]
            }
        },
        watch: {
            search(val) {
                if (val && val.length > 0) this.getList()
            }
        },
        methods: {
            getList() {
                postRequest({query: this.search}).subscribe(res => {
                    if (res) {
                        if (!res) res = []
                        this.resultList = res
                    }
                })
            },
            openDialogWithDetails(item) {
                this.selectedItem = item
                this.isShowDialog = true
            }
        },
    }
    const postRequest = ({query}) => ajax({
            url: `https://suggestions.dadata.ru/suggestions/api/4_1/rs/findById/party`,
            method: 'POST',
            headers: getHttpHeaders(),
            body: {
                query,
            }
        }).pipe(
            take(1),
            map(processResponse()),
            catchError(processError())
        )

    const getHttpHeaders = () => {
        let headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Token ${config.dadataToken}`,
        }
        return headers
    }

    const processResponse = ({successMsg} = {}) => (res) => {
        if (!res.response) throw new Error(res.response.message)
        return res.response.suggestions
    }

    const processError = (isShowError) => (err) => {
        const message = err.response ? err.response.message : err.message
        if (isShowError) {
            Notify.create({
                color: 'negative',
                position: 'top-right',
                message
            })
        }
        return of({ok: false, message})
    }
</script>
