<template>
    <div>
        <q-bar class="bg-secondary text-white shadow-2">
            <div>{{label}}</div>
            <q-space />
            <q-btn v-if = '!isDeleted' dense round flat icon="delete" @click="isDeleted = !isDeleted"><q-tooltip>показать список удаленных</q-tooltip></q-btn>
            <q-btn v-else dense round outline icon="delete" @click="isDeleted = !isDeleted"><q-tooltip>показать список активных</q-tooltip></q-btn>
            <q-btn v-if="!readonly" icon="add" round flat @click="add"><q-tooltip>Добавить</q-tooltip></q-btn>
        </q-bar>

        <q-list bordered separator dense>
            <q-item v-for="item in filteredList" :key="item.id">
                [[if GetJsonList.Icon]]
                <q-item-section avatar>
                    <q-avatar rounded>
                        <img src="[[GetJsonList.Icon]]">
                    </q-avatar>
                </q-item-section>
                [[end]]

                <!--  поля формы    -->
                <q-item-section class="col-[[if GetJsonList.Icon]]10[[else]]11[[end]]">
                    <div class="row q-col-gutter-md">
                    [[- range GetJsonList.Flds]]
                        [[PrintVueFldTemplate .]]
                    [[- end]]
                    </div>
                </q-item-section >

                <q-item-section v-if="!readonly" side>
                    <q-btn icon="delete" size="sm" v-if="!item.deleted" round flat @click="removeRecover(item.id)"><q-tooltip>{{$t('message.delete')}}</q-tooltip></q-btn>
                    <q-btn icon="check_circle_outline" size="sm" v-if="item.deleted" round flat @click="removeRecover(item.id)"><q-tooltip>Восстановить</q-tooltip></q-btn>
                </q-item-section>
            </q-item>
        </q-list>

    </div>

</template>

<script>
    import _ from 'lodash'

    export default {
        props: ['item', 'fld', 'label', 'readonly'],
        computed: {
            filteredList() {
                return this.list ? this.list.filter(v => v.deleted === this.isDeleted) : []
            }
        },
        data() {
            return {
                list: null,
                isDeleted: false,
            }
        },
        watch: {
            list: {
                handler() {
                    this.$emit('update', this.list)
                },
                deep: true,
            },
        },
        methods: {
            add() {
                let id = this.list.length > 0 ? _.maxBy(this.list, 'id').id : 0
                this.list.unshift({id: ++id, [[range GetJsonList.Flds]] [[.Name]]: null,[[end -]] deleted: false })
            },
            removeRecover(id) {
                let item = this.list.find(v => v.id === id)
                if (item) item.deleted = !item.deleted
            },
        },
        mounted() {
            this.list = this.fld || []
        }
    }
</script>
