<template>
  <q-page :padding="!isOpenInDialog">
    <comp-breadcrumb v-if="!isOpenInDialog" :list="[{label:'[[index .Vue.I18n "listTitle"]]', docType:'[[.Name]]'}]"/>

    <comp-doc-list ref="docList" listTitle='[[index .Vue.I18n "listTitle"]]' listDeletedTitle='[[index .Vue.I18n "listDeletedTitle"]]' pg-method="[[.PgName]]_list"
                   :list-sort-data="listSortData" :list-filter-data="listFilterData"
                   [[if not .Vue.IsHideCreateNewBtn]] :newDocUrl="currentUrl + 'new'" [[end]]
                   [[- if .IsRecursion]] :ext="ext ? Object.assign(ext, {parent_id: 'null'}) : {parent_id: 'null'}" [[else]] :ext="ext" [[end]]
                   search-fld-name="search_text">

      <template #listItem="{item}">
        [[.PrintListRowAvatar]]
        [[.PrintListRowLabel]]
        <q-item-section top side>
          <comp-item-dropdown-btn :item="item" itemProp="title" :is-edit="true" :is-delete="[[not .Vue.IsHideDeleteOptions]]" fkProp=""
                                  pg-method="[[.PgName]]_update"
                                  @edit="$router.push(`${currentUrl}${item.id}`)"
                                  @reload-list="$refs.docList.reloadList()"/>
        </q-item-section>
      </template>

    </comp-doc-list>
  </q-page>
</template>

<script>
  import currentUserMixin from '../../../app/mixins/currentUser'
  export default {
    props: ['isOpenInDialog', 'ext'],
    mixins: [currentUserMixin],
    computed: {
      currentUrl: () => '[[.Vue.RouteName]]/',
    },
    data() {
      return {
        listSortData: [
          {value: 'created_at', title: 'Дата'},
          {value: 'title', title: 'Название'}
        ],
        listFilterData: [
          {value: {deleted: false}, title: 'Активные'},
          {value: {deleted: true}, title: 'Удаленные'}
        ],
      }
    },
  }
</script>
