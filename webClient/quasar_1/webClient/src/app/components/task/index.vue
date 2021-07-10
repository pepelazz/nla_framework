<template>
  <q-page padding>
    <comp-breadcrumb :list="[{label:'Задача', docType: 'task', icon: 'fas fa-tasks'}]"/>

    <comp-doc-list ref="docList" doc-name="task" pg-method="task_list"
                   :list-sort-data="listSortData" :list-filter-data="listFilterData"
                   :newDocUrl="currentUrl + 'new'"
                   search-fld-name="search_text">

      <template #listItem="{item}">
        <q-item-section avatar @click="$router.push(`${currentUrl}${item.id}`)">
            <q-avatar rounded>
              <img src="https://image.flaticon.com/icons/svg/1642/1642808.svg">
            </q-avatar>
        </q-item-section>

        <q-item-section>
          <q-item-label lines="1">{{item.task_type_title}}</q-item-label>
          <q-item-label v-if="item.executor_fullname" caption><q-icon name="person"/> {{item.executor_fullname}}</q-item-label>
<!--          <q-item-label v-if="item.table_name == 'client'" caption @click="$router.push(`/client/${item.table_id}`)"><q-icon name="far fa-building"/> {{item.table_options.title}}</q-item-label>-->
        </q-item-section>

        <q-item-section top side>
          <comp-item-dropdown-btn :item="item" itemProp="title" :is-edit="true" :is-delete="true" fkProp=""
                                  pg-method="task_update"
                                  @edit="$router.push(`${currentUrl}${item.id}`)"
                                  @reload-list="$refs.docList.reloadList()"/>
        </q-item-section>
      </template>

    </comp-doc-list>
  </q-page>
</template>

<script>
  export default {
    computed: {
      currentUrl: () => '/task/',
    },
    data() {
      return {
        listSortData: [
          {value: 'created_at', title: 'Дата'},
          {value: 'Task', title: 'Название'}
        ],
        listFilterData: [
          {value: {deleted: false}, title: 'Активные'},
          {value: {deleted: true}, title: 'Удаленные'}
        ],
      }
    },
  }
</script>
