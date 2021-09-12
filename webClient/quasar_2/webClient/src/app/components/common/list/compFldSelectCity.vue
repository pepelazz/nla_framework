<template>
  <div>
    <q-select
      outlined
      v-model="localItem"
      use-input
      hide-selected
      fill-input
      input-debounce="300"
      :label='label'
      :options='options'
      @filter="filterFn"
      :readonly="readonly"
    >
      <!-- кнопка добавления города   -->
      <template v-slot:append v-if="!(ext && ext.isHideAddBtn)">
        <q-btn round dense flat icon="add" @click="isShowDialog=true"/>
      </template>

      <template v-slot:no-option>
        <q-item>
          <q-item-section class="text-grey">
            Ничего не найдено
          </q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- диалог добавления   -->
    <q-dialog v-model="isShowDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <q-input outlined label="Название города" v-model="newCityName"/>
        </q-card-section>
        <q-card-section>
        </q-card-section>
        <q-card-actions align="right" class="text-primary">
          <q-btn flat :label="$t('message.cancel')" v-close-popup/>
          <q-btn flat label="Добавить" v-close-popup @click="add"/>
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
    export default {
        props: {
            item: {
                type: String
            },
            ajaxSelectTitle: {
                type: String
            },
            label: {
                type: String
            },
            ext: {
                type: Object,
            },
            readonly: {
                type: Boolean,
                default: false
            },
        },
        data() {
            return {
                localItem: {},
                options: [],
                isShowDialog: false,
                newCityName: null,
            }
        },
        created() {
            const title = this.ajaxSelectTitle ? this.ajaxSelectTitle : null
            this.localItem = {label: title, value: title}
        },
        watch: {
            localItem: function (v) {
                if (v && v.id) this.$emit('update', v.id)
            }
        },
        methods: {
            filterFn(val, update, abort) {
                update(() => {
                    this.$utils.postCallPgMethod({
                        method: 'city_list',
                        params: Object.assign({search_text: val, per_page: 20}, this.ext ? this.ext : {}),
                    }).subscribe(res => {
                        if (res.ok) {
                            if (!res.result) res.result = []
                            this.options = res.result.map(v => {
                                return {
                                    label: v.title,
                                    value: v.title,
                                    id: v.id,
                                }
                            })
                        }
                    })
                })
            },
            add() {
                this.$utils.postCallPgMethod({method: 'city_update', params: {id: -1, title: this.newCityName}, isShowError: false}).subscribe(res => {
                    if (res.ok) {
                        const city = res.result
                        this.localItem = {label: city.title, value: city.id}
                        this.newCityName = null
                    } else {
                        if (res.message.includes('city_already_exist')) res.message = 'Город с таким назваанием уже существует'
                        this.$q.notify({
                            color: 'negative',
                            position: 'bottom',
                            message: res.message,
                        })
                    }
                })
            }
        },
    }
</script>
