<template>
  <q-page padding>
    <comp-breadcrumb :list="[{label:'Пользователи', docType: 'users'}]"/>

    <comp-doc-list ref="docList" doc-name="user" pg-method="user_list"
                   :list-sort-data="listSortData" :list-filter-data="listFilterData"
                   search-fld-name="search_fullname">

      <template #listItem="{item}">
        <q-item-section avatar @click="$router.push(`${currentUrl}${item.id}`)">
          <q-avatar rounded>
            <comp-stat-img-src v-if="item.avatar" :src="item.avatar"/>
            <img v-else src="https://www.svgrepo.com/show/95333/avatar.svg">
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label lines="1">{{item.fullname}}</q-item-label>
          <q-item-label caption lines="2">{{item.email}}</q-item-label>
          <q-item-label caption lines="2">{{rolesI18n(item)}}</q-item-label>
        </q-item-section>
        <q-item-section top side>
          <comp-item-dropdown-btn :item="item" itemProp="fullname" :is-edit="true" :is-delete="true"
                                  pg-method="user_update"
                                  @edit="$router.push(`${currentUrl}${item.id}`)"
                                  @reload-list="$refs.docList.reloadList()"/>
        </q-item-section>
      </template>

    </comp-doc-list>
  </q-page>
</template>

<script>
  import roles from './roles'
  import _ from 'lodash'
  export default {
    computed: {
      currentUrl: () => '/users/',
    },
    data() {
      return {
        listSortData: [
          {value: 'created_at', title: 'Дата'},
          {value: 'fullname', title: 'ФИО'}
        ],
        listFilterData: [
          {value: {deleted: false}, title: 'Активные'},
          {value: {deleted: true}, title: 'Удаленные'}
        ],
      }
    },
    methods: {
        rolesI18n(item) {
            return item.role.filter(v => v !== 'student').map(v => _.find(roles, {value: v}).label)
        }
    },
  }
</script>
