<template>
    <div>
        <q-bar class="bg-secondary text-white shadow-2">
            <div>[[GetLabel]] <span v-if="deleted">удаленные</span></div>
            <q-space />
            [[if IsShowAdd]]<q-btn icon="add" v-if="!readonly" round flat @click="add"><q-tooltip>Добавить</q-tooltip></q-btn>[[- end]]
            [[if IsShowDelete]]<q-btn v-if="deleted && !readonly" icon="delete" round flat @click="reload(false)"><q-tooltip>активные [[GetLabel]]</q-tooltip></q-btn>[[- end]]
            [[if IsShowDelete]]<q-btn v-if="!deleted && !readonly" icon="delete_outline" round flat @click="reload(true)"><q-tooltip>удаленные [[GetLabel]]</q-tooltip></q-btn>[[- end]]
        </q-bar>

        <q-list bordered separator>
            <q-item v-for="v in list" :key="v.id">
                <q-item-section avatar @click="$router.push('/[[GetRoute]]/' + v.id)" style="cursor: pointer">
                    <q-avatar rounded size="sm">
                        <img src="[[GetAvatar]]" alt="">
                    </q-avatar>
                </q-item-section>
                <q-item-section>
                    [[GetTitleTemplate]]
                </q-item-section>
                [[if IsShowDelete]] <q-item-section side v-if="!readonly">
                    <q-icon :name="deleted ? 'done' : 'delete'" size="xs" class="cursor-pointer" color="grey" @click="removeRecover(v)"/>
                </q-item-section>[[end]]
            </q-item>
        </q-list>

        <!-- диалог добавления       -->
        [[if IsShowAdd]]<q-dialog v-model="isShowAddDialog">
            <q-card style="width: 500px; max-width: 80vw;">
                <q-bar>
                    <div>Создать новую запись</div>
                    <q-space />
                    <q-btn dense flat icon="close" v-close-popup/>
                </q-bar>
                <q-card-section>
                    [[range GetNewFlds]]
                    [[PrintVueFldTemplate .]]
                    [[end]]
                </q-card-section>
                <q-card-actions align="right" class="bg-white text-teal">
                    <q-btn flat label="OK" @click="saveNew"/>
                </q-card-actions>
            </q-card>
        </q-dialog>[[- end]]
    </div>
</template>

<script>
    [[range GetTagFlds]]
    import [[.]]_tag_list from '../../[[GetTableName]]/mixins/[[.]]_tag_list'
    [[- end]]
    export default {
        props: ['id', 'readonly'],
        mixins: [ [[range GetTagFlds]][[.]]_tag_list,[[end]] ],
        data() {
            return {
                list: [],
                isShowAddDialog: false,
                deleted: false,
                item: {[[range GetNewFlds]][[.Name]]: null, [[end]]},
            }
        },
        methods: {
            add() {
                this.isShowAddDialog = true
            },
            reload(isDeleted) {
                !isDeleted ? this.deleted = false : this.deleted = true
                this.$utils.callPgMethod('[[GetTableName]]_list', {'[[GetRefFldName]]': this.id, deleted: this.deleted, 'order_by': 'created_at desc'}, (result) => this.list = result)
            },
            saveNew() {
                [[range GetNewFlds]]
                    [[if .Vue.IsRequired]]if (!this.item.[[.Name]]) {
                    this.$q.notify({type: 'negative', message: 'не заполнено поле: "[[.NameRu]]"'})
                    return
                }[[- end]]
                [[- end]]
                let params = Object.assign({id: -1, [[GetRefFldName]]: this.id}, this.item)
                [[range GetNewFlds]]
                    [[if eq .Vue.Type "select"]]
                        // если поле select то заменяем объект на строку
                        if (this.item.[[.Name]]) params = Object.assign(params, {[[.Name]]:this.item.[[.Name]].value})
                        [[- end]]
                [[- end]]
                // если IsStateMachine то [[GetTableName]]_create, в остальных случаях [[GetTableName]]_update
                this.$utils.callPgMethod('[[GetTableName]]_[[if IsStateMachine]]create[[else]]update[[end]]', params, () => {
                    this.isShowAddDialog = false
                    [[range GetNewFlds]]
                    this.item.[[.Name]] = null [[end]]
                    this.reload()
                })
            },
            removeRecover({id}) {
                this.$utils.callPgMethod('[[GetTableName]]_update', {id, deleted: !this.deleted}, () => this.reload(this.deleted))
            }
        },
        mounted() {
            this.reload()
        }
    }
</script>
