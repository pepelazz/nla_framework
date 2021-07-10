<template>
  <div>
    <q-input
      dense
      :label='label'
      v-model="search"
      filled
      debounce="300"
    >
      <template v-slot:prepend v-if="isFirstSelected">
        <q-icon name="add" @click="add" />
      </template>
      <template v-slot:append>
        <q-icon name="search" />
      </template>
    </q-input>
    <p class="text-caption q-mb-none" v-if="notFound">не найдено</p>
    <p class="text-caption q-mb-none" v-if="resultList.length>1">выберите вариант из списка или продолжите ввод</p>
    <q-list bordered separator v-if="resultList.length>0">
      <q-item v-for="(item) in resultList" :key="item.data.hid" clickable v-ripple>
        <q-item-section @click="selectAddress(item)">
          <q-item-label>{{item.value}}</q-item-label>
          <!--          <q-item-label caption>{{item.data.inn}} {{item.data.address.value}}</q-item-label>-->
        </q-item-section>
      </q-item>
    </q-list>
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
                search: '',
                resultList: [],
                selectedItem: null,
                notFound: false,
                isFirstSelected: false, // флаг для выделения, что как минимум один выборо из списка уже сделан
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
                        this.notFound = (res.length === 0)
                        if (res.length === 0) this.isFirstSelected = false
                    }
                })
            },
            selectAddress(item) {
                this.search = item.value
                this.isFirstSelected = true
                this.$nextTick(() => {
                    if (this.resultList.length === 1) this.resultList = []
                })
            },
            // специальный запрос чтобы получить геокоординаты
            add() {
                postRequest({query: this.search, count: 1}).subscribe(res => {
                    if (res) {
                        this.resultList = []
                        this.search = ''
                        this.$emit('update', res)
                    }
                })
            }
        },
    }
    const postRequest = ({query, count = 10}) => ajax({
        url: `https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address`,
        method: 'POST',
        headers: getHttpHeaders(),
        body: {
            query,
            count
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
