<template>
    <div  v-if="id != 'new'">
        <div class="row q-col-gutter-md q-mb-sm">
            <div class="col">
                <q-btn icon="add" round flat @click="$refs.addTaskDialog.open()">
                    <q-tooltip>Добавить задачу</q-tooltip>
                </q-btn>
            </div>
        </div>
        <!--  поля формы    -->
        <div class="row q-col-gutter-md q-mb-sm">
            <div class="col-md-6">
                <q-list bordered separator>
                    <q-item v-for="item in list" :key="item.id">
                        <q-item-section avatar top @click="$router.push(`/task/${item.id}`)">
                            <q-avatar v-if="item.isDeadlinePass" icon="warning" color="orange" text-color="white"/>
                            <q-avatar v-else icon="error_outline" color="info" text-color="white"/>
                        </q-item-section>
                        <q-item-section>
                            <q-item-label>{{item.task_type_title}}</q-item-label>
                            <q-item-label caption>Исполнитель: {{item.executor_fullname}}</q-item-label>
                            <q-item-label caption>Deadline: <strong>{{$utils.formatPgDateTime(item.deadline)}}</strong></q-item-label>
                        </q-item-section>
                        <q-item-section side>
                            <div class="text-grey-8 q-gutter-xs">
                                <q-btn class="gt-xs" size="12px" flat outline round icon="check" @click="$refs.doneTaskDialog.open(item)"/>
                                <!--                  <q-btn class="gt-xs" size="12px" flat dense round icon="edit" @click="$router.push(`/task/${item.id}`)"/>-->
                            </div>
                        </q-item-section>
                    </q-item>
                </q-list>
                <!--          <pre>{{list}}</pre>-->
            </div>
        </div>
        <comp-dialog-task-done ref="doneTaskDialog" @taskFinished="v=>$emit('taskFinished', v)"/>
        <comp-dialog-task-add ref="addTaskDialog" :table_id="id" table_name="[[.Name]]"/>
    </div>
</template>

<script>
    export default {
        props: ['id', 'list'],
        computed: {
            docUrl: () => '/[[.Vue.RouteName]]',
        },
        data() {
            return {
                tab: 'tasks',
            }
        },
        methods: {
        },
    }
</script>
