<template>
    <q-item>
      <q-item-section avatar top @click="$router.push(`/task/${item.id}`)">
        <template v-if="item && item.task_type_options && item.task_type_options.iconUrl">
          <q-avatar rounded>
            <q-img :src="item.task_type_options.iconUrl"/>
          </q-avatar>
        </template>
        <template v-else>
          <q-avatar v-if="item.isDeadlinePass" rounded icon="warning" color="orange" text-color="white"/>
          <q-avatar v-else rounded icon="error_outline" color="info" text-color="white"/>
        </template>
      </q-item-section>
      <q-item-section>
        <q-item-label>{{item.task_type_title}}</q-item-label>
        <q-item-label caption>{{$utils.formatPgDate(item.deadline)}}</q-item-label>
<!--        <q-item-label v-if="item.table_name && item.table_options" caption @click="$router.push(`/${item.table_name}/${item.table_id}`)"><q-icon :name="icon(item.table_name)"/> {{item.table_options.title}}</q-item-label>-->
        <!--            <q-item-label v-if="item.table_name === 'client'" caption @click="$router.push(`/client/${item.table_id}`)"><q-icon name="far fa-building"/> {{item.table_options.title}}</q-item-label>-->
        <!--            <q-item-label v-if="item.table_name === 'deal'" caption @click="$router.push(`/client/${item.table_options.client_id}/deal/${item.table_id}`)"><q-icon name="opacity"/> {{item.table_options.client_title}} {{item.table_options.deal_title}}</q-item-label>-->
      </q-item-section>
      <q-item-section side v-if="item.state != 'finished'">
        <div class="text-grey-8">
          <q-btn size="12px" flat dense round icon="done" @click="$refs.doneTaskDialog.open(item)"/>
        </div>
      </q-item-section>
      <comp-dialog-task-done ref="doneTaskDialog" @taskFinished="v=>$emit('taskFinished', v)"/>
    </q-item>
</template>

<script>
  export default {
      props: ['item'],
      data() {
          return {
          }
      },
      methods: {
      },
  }
</script>
