<template>
  <q-page padding>
    <comp-breadcrumb :list="[{label:'Типы задач', docType: 'taskType', icon: 'bookmark'}]"/>

    <comp-doc-list ref="docList" doc-name="taskType" pg-method="task_type_list"
                   :list-sort-data="listSortData" :list-filter-data="listFilterData"
                   :newDocUrl="currentUrl + 'new'"
                   search-fld-name="search_text">

      <template #listItem="{item}">
        <q-item-section avatar @click="$router.push(`${currentUrl}${item.id}`)">
          <q-avatar rounded>
            <img v-if="item.options && item.options.iconUrl" :src="item.options.iconUrl">
            <img v-else src="https://image.flaticon.com/icons/svg/1030/1030881.svg">
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label lines="1">{{item.title}}</q-item-label>
          <q-item-label v-if="item.table_name" caption>{{$config.tablesForTask[item.table_name]}}</q-item-label>

        </q-item-section>
        <q-item-section top side>
          <comp-item-dropdown-btn :item="item" itemProp="title" :is-edit="true" :is-delete="true" fkProp=""
                                  pg-method="task_type_update"
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
      currentUrl: () => '/taskType/',
    },
    data() {
      return {
        listSortData: [
          {value: 'created_at', title: 'Дата'},
          {value: 'TaskType', title: 'Название'}
        ],
        listFilterData: [
          {value: {deleted: false}, title: 'Активные'},
          {value: {deleted: true}, title: 'Удаленные'}
        ],
      }
    },
  }
</script>
