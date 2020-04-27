<template>
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
    <!-- кнопка добавления   -->
    <template v-slot:append v-if="ext && ext.addNewUrl">
      <q-btn round dense flat icon="add" @click="openNewTab"/>
    </template>

    <template v-slot:no-option>
      <q-item>
        <q-item-section class="text-grey">
          Ничего не найдено
        </q-item-section>
      </q-item>
    </template>
  </q-select>
</template>

<script>
    export default {
        props: {
            item: {
                type: String
            },
            label: {
                type: String
            },
            ext: {
                type: Object
            },
            readonly: {
                type: Boolean,
                default: false
            },
            itemTitleFldName: {
                type: String,
                default: 'title'
            },
            pgMethod: {
                type: String
            },
        },
        data() {
            return {
                localItem: {},
                options: []
            }
        },
        created() {
            const title = this.item ? this.item : null
            this.localItem = {label: title, value: title}
        },
        watch: {
            localItem: function (v) {
                if (v && v.id) this.$emit('update', {id: v.id, item: v})
            }
        },
        methods: {
            filterFn(val, update, abort) {
                update(() => {
                    this.$utils.postCallPgMethod({
                        method: this.pgMethod,
                        params: Object.assign({search_text: val, per_page: 20}, this.ext ? this.ext : {}),
                    }).subscribe(res => {
                        if (res.ok) {
                            if (!res.result) res.result = []
                            this.options = res.result.map(v => {
                                return {
                                    label: v[this.itemTitleFldName],
                                    value: v[this.itemTitleFldName],
                                    id: v.id,
                                }
                            })
                        }
                    })
                })
            },
            openNewTab() {
                window.open(this.ext.addNewUrl, '_blank')
            }
        },
    }
</script>
