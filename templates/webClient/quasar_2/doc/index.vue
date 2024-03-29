[[$doc:=.]]
[[$parent:=.]]


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
        [[- else if .IsDate]]
        <div class="[[if .ColClass]] [[.ColClass]] [[else]] col-md-2 col-sm-4 col-xs-6 [[- end]]">
          <comp-fld-date outlined label="[[.Label]]" :date-string="$utils.formatPgDate(filter[[ToCamel .FldName]])" @update="updateFilter[[ToCamel .FldName]]" @clear="updateFilter[[ToCamel .FldName]](null)" dense/>
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
                   [[if .Vue.CreateNewModal.Flds]] :new-doc-event-only="true" @clickAddBtn="clickAddBtn" [[- end]]
                  [[if and (not .Vue.IsHideCreateNewBtn) (not .Vue.CreateNewModal.Flds)]] :newDocUrl="currentUrl + 'new'" [[- end]]
                  [[if .Vue.IsOpenNewInTab]] :isOpenNewInTab="true" [[- end]]
                   [[- if .Vue.ListUrlQueryParams]] :urlQueryParams="[ [[range .Vue.ListUrlQueryParams]]'[[.]]',[[- end]] ]" [[end]]
                   [[- if .IsRecursion]] :ext="ext ? Object.assign(ext, {parent_id: 'null'}) : {parent_id: 'null'}" [[else]] :ext="ext" [[end]]
                   [[- if .Vue.List.ColClass]] col-class="[[.Vue.List.ColClass]]" [[- end]]
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

  [[- if .Vue.CreateNewModal.Flds]]
  <!-- диалог добавления -->
  <comp-edit-dialog ref="addDialogRef" labelNew="[[- if .Vue.CreateNewModal.Label]] [[- .Vue.CreateNewModal.Label]] [[- else]]Новая [[.NameRu]] [[- end]]" pgMethod="[[if .Vue.CreateNewModal.PgMethod]] [[- .Vue.CreateNewModal.PgMethod]] [[- else -]] [[.Name]]_update [[- end]]" @update="v => $router.push('/[[.Name]]/' + v.id)"/>
  [[- end]]

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
        [[$parent = .Vue]]
        listSortData: [
          [[- range .Vue.SortList]]
          {value: '[[.Value]]', title: [[$parent.PrintSortListLabel .Label]]},
          [[- end]]
        ],
        [[- end]]
        listFilterData: [
          {value: {deleted: false}, title: this.$t('message.filter_active')},
          {value: {deleted: true}, title: this.$t('message.filter_deleted')}
        ],
        [[- range .Vue.FilterList]]
        [[- if .IsRef]]
        filter[[ToCamel .RefTable]]Title: null,
        filter[[ToCamel .RefTable]]Id: null,
        [[- else if .IsDate]]
        filter[[ToCamel .FldName]]: null,
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
        [[- if .IsSaveLocalStorage]]
        localStorageKey_[[ToCamel $doc.Name]][[ToCamel .FldName]]: 'filter_[[ToCamel $doc.Name]][[ToCamel .FldName]]', [[- end]]
        [[- end]]
      }
    },
    methods: {
      [[- range .Vue.FilterList]]
      [[- if .IsRef]]
      updateFilter[[ToCamel .RefTable]](v) {
        this.$refs.docList.changeItemList({'[[.FldName]]': v ? v.id : null})
        [[- if .IsSaveLocalStorage]]
        localStorage.setItem(this.localStorageKey_[[ToCamel $doc.Name]][[ToCamel .FldName]], v ? v.id : null)
        [[- end]]
        if (v) {
          this.$utils.callPgMethod(`[[.RefTable]]_get_by_id`, {id: v.id}, (res) => {
            this.filter[[ToCamel .RefTable]]Title = res.title
          })
        }
      },
      [[- else if .IsDate]]
      updateFilter[[ToCamel .FldName]](v) {
        this.$refs.docList.changeItemList({'[[.FldName]]': v ? v : null})
        this.filter[[ToCamel .FldName]] = v
      },
      [[- else]]
      updateFilter[[ToCamel .FldName]](v) {
        this.$refs.docList.changeItemList({'[[.FldName]]': v ? v : null})
        this.filter[[ToCamel .FldName]] = v ? this.options[[ToCamel .FldName]].find(v1 => v1.value === v) : null
      },
      [[- end]]
      [[- end]]

      [[- if .Vue.CreateNewModal.Flds]]
      clickAddBtn() {
        this.$refs.addDialogRef.show({
          item: {id: -1, title: null, to_user_id: null},
          flds: [
            [[- range .Vue.CreateNewModal.Flds]]
            [
              [[- range .]]
              [[$doc.PrintCompEditDialogFld .]],
              [[- end]]
            ],
            [[- end]]
          ],
        })
      }
      [[- end]]
    },
    mounted() {
    [[if .Vue.FilterList -]]
      // извлекаем параметры фильтрации из url
      const urlParams = new URLSearchParams(window.location.search)
      [[- range .Vue.FilterList]]
      [[- if .IsRef]]
        [[- if .IsSaveLocalStorage]]
          const [[ToCamel $doc.Name]][[ToCamel .FldName]]Id = localStorage.getItem(this.localStorageKey_[[ToCamel $doc.Name]][[ToCamel .FldName]])
          if (+[[ToCamel $doc.Name]][[ToCamel .FldName]]Id) {
            this.updateFilter[[ToCamel .RefTable]]({id: +[[ToCamel $doc.Name]][[ToCamel .FldName]]Id})
          }
        [[- else]]
          if (urlParams.has('[[.FldName]]')) {
            let id = +urlParams.get('[[.FldName]]')
            if (id) this.updateFilter[[ToCamel .RefTable]]({id})
          }
        [[- end]]
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
