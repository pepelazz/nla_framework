<template>
  <q-page :padding="!isOpenInDialog">
    <comp-breadcrumb v-if="!isOpenInDialog" :list="[{label:'[[index .Vue.I18n "listTitle"]]', docType:'[[.Name]]'}]"/>

    [[- if .Vue.FilterList]]
    <!-- фильтры   -->
    <div class="row q-mt-sm q-col-gutter-sm">
      [[- range .Vue.FilterList]]
        [[- if .IsRef]]
        <div class="[[.ColClass]]">
          <comp-fld-ref-search dense outlined pgMethod="[[.RefTable]]_list" label="[[.Label]]" :item='filter[[ToCamel .RefTable]]Title' :itemId='filter[[ToCamel .RefTable]]Id' :ext='{isClearable: true}'  @update="updateFilter[[ToCamel .RefTable]]" @clear="updateFilter[[ToCamel .RefTable]]"  class='q-mb-sm col-md-4 col-sm-6 col-xs-12' />
        </div>
        [[- end]]
      [[- end]]
    </div>
    [[- end]]

    <comp-doc-list ref="docList" listTitle='[[index .Vue.I18n "listTitle"]]' listDeletedTitle='[[index .Vue.I18n "listDeletedTitle"]]' pg-method="[[.PgName]]_list"
                   :list-sort-data="listSortData" :list-filter-data="listFilterData"
                  [[if not .Vue.IsHideCreateNewBtn]] :newDocUrl="currentUrl + 'new'" [[- end]]
                  [[if .Vue.IsOpenNewInTab]] :isOpenNewInTab="true" [[- end]]
                   [[- if .Vue.ListUrlQueryParams]] :urlQueryParams="[ [[range .Vue.ListUrlQueryParams]]'[[.]]',[[- end]] ]" [[end]]
                   [[- if .IsRecursion]] :ext="ext ? Object.assign(ext, {parent_id: 'null'}) : {parent_id: 'null'}" [[else]] :ext="ext" [[end]]
                   search-fld-name="search_text" :readonly="[[.Vue.Readonly]]">

      [[- if .Vue.List.AddBtnsSlot]]
      <template #addBtnsSlot>
        [[range .Vue.List.AddBtnsSlot]]
        [[- if .UploadFile.Url]]<comp-file-upload url='[[.UploadFile.Url]]' :file-ext='[ [[ArrayStringJoin .UploadFile.FileExt]] ]' tooltip="[[.UploadFile.Tooltip]]" style="display: contents" @reloadList="$refs.docList.reloadList()"/>[[- end]]
        [[- end]]
      </template>
      [[- end]]


      <template #listItem="{item}">
        [[.PrintListRowAvatar]]
        [[.PrintListRowLabel]]
        <q-item-section top side>
          <comp-item-dropdown-btn :item="item" itemProp="title" :is-edit="true" :is-delete="!([[.Vue.Readonly]] || [[.Vue.IsHideCreateNewBtn]])" fkProp=""
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
      currentUrl: () => '/[[.Vue.RouteName]]/',
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
        [[- range .Vue.FilterList]]
        [[- if .IsRef]]
        filter[[ToCamel .RefTable]]Title: null,
        filter[[ToCamel .RefTable]]Id: null,
        [[- end]]
        [[- end]]
      }
    },
    methods: {
      [[- range .Vue.FilterList]]
      [[- if .IsRef]]
      updateFilter[[ToCamel .RefTable]](v) {
        this.$refs.docList.changeItemList({'[[.FldName]]': v ? v.id : null})
        if (v) {
          this.$utils.callPgMethod(`[[.RefTable]]_get_by_id`, {id: v.id}, (res) => {
            this.filter[[ToCamel .RefTable]]Title = res.title
          })
        }
      },
      [[- end]]
      [[- end]]
    },
    mounted() {
    [[if .Vue.FilterList -]]
      // извлекаем параметры фильтрации из url
      const urlParams = new URLSearchParams(window.location.search)
      [[- range .Vue.FilterList]]
      [[- if .IsRef]]
      if (urlParams.has('[[.FldName]]')) {
        let id = +urlParams.get('[[.FldName]]')
        if (id) this.updateFilter[[ToCamel .RefTable]]({id})
      }
      [[- end]]
    [[- end]]
    [[- end]]
    }
  }
</script>
