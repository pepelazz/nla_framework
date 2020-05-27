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

    <!-- аватарка со ссылкой   -->
    <template v-slot:before v-if="ext && ext.pathUrl">
      <q-avatar rounded v-if="(localItem && localItem.id) || itemId" @click="$router.push(ext.pathUrl + '/' + ((localItem && localItem.id) || itemId))">
        <img :src="ext.avatar">
      </q-avatar>
      <q-avatar v-else rounded style="opacity: 0.7"><img :src="ext.avatar"></q-avatar>
    </template>

    <!-- кнопка добавления   -->
    <template v-slot:append v-if="ext && ext.addNewUrl">
      <q-btn round dense flat icon="add" @click="openNewTab"/>
    </template>

    <!-- кнопка удаления   -->
    <template v-slot:append  v-if="ext && ext.isClearable && localItem.label">
      <q-icon name="cancel" @click.stop="clear" class="cursor-pointer" />
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
            itemId: {
                type: Number
            },
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
            clear() {
                this.localItem = {label: null, value: null}
                this.$emit('clear')
            },
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
                                    item: v,
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
        mounted() {
            // подгружаем список сразу при открытии
            this.$utils.postCallPgMethod({
                method: this.pgMethod,
                params: Object.assign({search_text: '', per_page: 20}, this.ext ? this.ext : {}),
            }).subscribe(res => {
                if (res.ok) {
                    if (!res.result) res.result = []
                    this.options = res.result.map(v => {
                        return {
                            label: v[this.itemTitleFldName],
                            value: v[this.itemTitleFldName],
                            id: v.id,
                            item: v,
                        }
                    })
                    // если результат в списке 1, то сразу ставим его как выбранный
                    // if (this.options.length === 1) {
                    //     this.localItem = this.options[0]
                    // }
                }
            })
        }
    }
</script>
