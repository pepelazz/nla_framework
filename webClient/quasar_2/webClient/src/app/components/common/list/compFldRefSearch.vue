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
    :dense="dense"
  >

    <!-- аватарка со ссылкой   -->
    <template v-slot:before v-if="ext && ext.pathUrl">
      <router-link v-if="(localItem && localItem.id) || itemId" :to="ext.pathUrl + '/' + ((localItem && localItem.id) || itemId)" style="cursor: pointer">
        <q-avatar rounded >
            <img :src="ext.avatar">
        </q-avatar>
      </router-link>
      <q-avatar v-else rounded style="opacity: 0.7"><img :src="ext.avatar"></q-avatar>
    </template>

    <template v-slot:append>
      <!-- кнопка добавления   -->
      <q-btn v-if="ext && ext.addNewUrl && !readonly" round dense flat icon="add" @click="openNewTab"/>
      <!-- кнопка удаления   -->
      <q-icon v-if="ext && ext.isClearable && localItem.label && !readonly" name="cancel" @click.stop="clear" class="cursor-pointer" />
    </template>

    <!-- кнопка сохранения   -->
    <template v-slot:after  v-if="ext && ext.isSaveBtn && localItem.label && !readonly">
      <q-icon name="save" @click="$emit('save')" class="cursor-pointer" color="primary" />
    </template>

    <!-- форматирование списка   -->
    <template v-slot:option="scope">
      <q-item
        v-bind="scope.itemProps"
        v-on="scope.itemEvents"
      >
        <q-item-section avatar v-if="scope.opt.icon">
          <q-icon :name="scope.opt.icon" />
        </q-item-section>
        <q-item-section>
          <q-item-label v-html="scope.opt.label" />
          <q-item-label caption v-if="ext && ext.descriptionFunc">{{ ext.descriptionFunc(scope.opt.item)}}</q-item-label>
        </q-item-section>
      </q-item>
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
            dense: {
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
            this.localItem = {label: title, value: this.ext?.value || title}
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
                                let label = v[this.itemTitleFldName]
                                if (this.ext?.itemTitleFunc) {
                                  label = this.ext.itemTitleFunc(v)
                                }
                                return {
                                    label,
                                    value: `${v[this.ext?.itemValueFldName ? this.ext.itemValueFldName : this.itemTitleFldName]}`,
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
                        let label = v[this.itemTitleFldName]
                        if (this.ext?.itemTitleFunc) {
                          label = this.ext.itemTitleFunc(v)
                        }
                        return {
                            label,
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
