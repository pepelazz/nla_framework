<template>
    <div  v-if="id != 'new'" class="row q-col-gutter-md q-mb-sm">
        <div class="col-6">
            <q-bar class="bg-secondary text-white shadow-2">
                <div>Задачи</div>
                <q-space />
                <q-btn dense flat round icon="fas fa-check-circle" v-if="selectedState === 'all'" @click="changeState('in_process')"/>
                <q-btn dense flat round icon="far fa-check-circle" v-if="selectedState === 'in_process'" @click="changeState('all')"/>
                <q-btn dense flat icon="add" @click="$refs.addTaskDialog.open()"/>
            </q-bar>
            <q-list separator bordered>
                <component v-for="item in listForRender" :key="item.id" :is="item.template" :item="item" @taskFinished="reload"></component>
            </q-list>
        </div>
        <comp-dialog-task-done ref="doneTaskDialog" @taskFinished="v=>$emit('taskFinished', v)"/>
        <comp-dialog-task-add ref="addTaskDialog" :table_id="id" table_name="[[.Name]]"/>
    </div>
</template>

<script>
    import defaultTmpl from '../../../currentUser/tasks/taskTemplates/defaultTmpl'
    export default {
        props: ['id'],
        components: {defaultTmpl},
        computed: {
            docUrl: () => '/[[.Vue.RouteName]]',
            listForRender: function () {
                return this.list.filter(v => this.selectedState === 'all' ? true : v.state === this.selectedState).map(v => {
                    v.template = v.options && v.options.template ? v.options.template : 'defaultTmpl'
                    return v
                })
            },
        },
        data() {
            return {
                list: [],
                tab: 'tasks',
                selectedState: 'in_process',
                tableName: '[[.Name]]',
            }
        },
        methods: {
            changeState(state) {
                this.selectedState = state
            },
            reload() {
                this.$utils.postCallPgMethod({method: 'task_list_by_[[.Name]]', params: {[[.Name]]_id: this.id}}).subscribe(res => {
                    if (res.ok) {
                        this.list = res.result.map(v => {
                            if (new Date(v.deadline) < new Date()) {
                                v.isDeadlinePass = true
                            }
                            return v
                        })
                    }
                })
            },
            formatDate(d) {
                return this.$utils.formatPgDate(d)
            },
        },
        mounted() {
            this.reload()
        }
    }
</script>
