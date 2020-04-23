<template>
  <q-page padding>
    <comp-breadcrumb :list="[{label:'Города'}]"/>

    <comp-doc-list ref="docList" doc-name="city" pg-method="city_list"
                   :list-sort-data="listSortData" :list-filter-data="listFilterData"
                   :newDocUrl="currentUrl + 'new'"
                   search-fld-name="search_text">

      <template #listItem="{item}">
        <q-item-section avatar @click="$router.push(`${currentUrl}${item.id}`)">
          <q-avatar rounded>
            <img src="https://image.flaticon.com/icons/svg/589/589554.svg" alt="">
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label lines="1">{{item.title}}</q-item-label>

        </q-item-section>
        <q-item-section top side>
          <comp-item-dropdown-btn :item="item" itemProp="" :is-edit="true" :is-delete="true" fkProp=""
                                  pg-method="city_update"
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
      currentUrl: () => '/city/',
    },
    data() {
      return {
        listSortData: [
          {value: 'created_at', title: 'Дата'},
          {value: 'city', title: 'Название'}
        ],
        listFilterData: [
          {value: {deleted: false}, title: 'Активные'},
          {value: {deleted: true}, title: 'Удаленные'}
        ],
      }
    },
  }
</script>
