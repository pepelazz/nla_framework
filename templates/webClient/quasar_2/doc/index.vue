[[$doc:=.]]

<template>
  <q-page :padding="!isOpenInDialog">
    <comp-breadcrumb class="text-capitalize" v-if="!isOpenInDialog" :list="[{label: $t('[[.Name]].name_plural'), docType:'[[.Name]]'}]"/>

    [[- if .Vue.FilterList]]
    <!-- фильтры   -->
    <div class="row q-mt-sm q-col-gutter-sm">
      [[- range .Vue.FilterList]]
        [[- if .IsRef]]
        <div class="[[if .ColClass]] [[.ColClass]] [[else]] col-md-2 col-sm-4 col-xs-6 [[- end]]">
          <comp-fld-ref-search dense outlined pgMethod="[[.RefTable]]_list" label="[[.Label]]" :item='filter[[ToCamel .RefTable]]Title' :itemId='filter[[ToCamel .RefTable]]Id' :ext='{isClearable: true}'  @update="updateFilter[[ToCamel .RefTable]]" @clear="updateFilter[[ToCamel .RefTable]]"  class='q-mb-sm col-md-4 col-sm-6 col-xs-12' />
        </div>
        [[- else]]
        <div class="[[if .ColClass]] [[.ColClass]] [[else]] col-md-2 col-sm-4 col-xs-6 [[- end]]">
          <q-select dense outlined v-model="filter[[ToCamel .FldName]]" :options="options[[ToCamel .FldName]]" label="[[.Label]]" @update:model-value="v => updateFilter[[ToCamel .FldName]](v.value)"  class='q-mb-sm col-md-4 col-sm-6 col-xs-12'>
            <template v-slot:append v-if="filter[[ToCamel .FldName]]">
              <q-icon name="close" @click.stop="updateFilter[[ToCamel .FldName]](null)" class="cursor-pointer" />
            </template>
          </q-select>
        </div>
        [[- end]]
      [[- end]]
    </div>
    [[- end]]

    <comp-doc-list ref="docList" :listTitle="$t('[[.Name]].name_plural')" :listDeletedTitle="$t('[[.Name]].name_plural_deleted')" pg-method="[[.PgName]]_list"
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
        [[- if .Comp.CompName]]<[[.Comp.CompName]] [[.Comp.Params]] style="display: contents" @reloadList="$refs.docList.reloadList()"/>[[- end]]
        [[- end]]
      </template>
      [[- end]]
      [[- if .Vue.List.AddFilterSlot]]
      <template #addFilterSlot>
        [[- range .Vue.List.AddFilterSlot]]
          [[.]]
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
[[ .PrintVueImport "docIndex" ]]
  import currentUserMixin from '../../../app/mixins/currentUser'
  export default {
    props: ['isOpenInDialog', 'ext'],
    components: {[[- .PrintComponents "docIndex" -]]},
    mixins: [currentUserMixin],
    computed: {
      currentUrl: () => '/[[.Vue.RouteName]]/',
    },
    data() {
      return {
        [[- if .Vue.SortList]]
        listSortData: [
          [[- range .Vue.SortList]]
          {value: '[[.Value]]', title: '[[.Label]]'},
          [[- end]]
        ],
        [[- end]]
        listFilterData: [
          {value: {deleted: false}, title: 'Активные'},
          {value: {deleted: true}, title: 'Удаленные'}
        ],
        [[- range .Vue.FilterList]]
        [[- if .IsRef]]
        filter[[ToCamel .RefTable]]Title: null,
        filter[[ToCamel .RefTable]]Id: null,
        [[- else]]
        filter[[ToCamel .FldName]]: null,
        options[[ToCamel .FldName]]: [
        [[- if .Options]]
          [[- range .Options]]
          {label: '[[.Label]]', value: '[[.Value]]'},
          [[- end]]
        [[- else]]
        [[ PrintFldSelectOptions $doc .FldName ]]
        [[- end]]
        ],
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
      [[- else]]
      updateFilter[[ToCamel .FldName]](v) {
        this.$refs.docList.changeItemList({'[[.FldName]]': v ? v : null})
        this.filter[[ToCamel .FldName]] = v ? this.options[[ToCamel .FldName]].find(v1 => v1.value === v) : null
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
      [[- else]]
      if (urlParams.has('[[.FldName]]')) {
        let name = urlParams.get('[[.FldName]]')
        if (name) this.updateFilter[[ToCamel .FldName]](name)
      }
      [[- end]]
    [[- end]]
    [[- end]]
    }
  }
</script>
