<template>
  <div v-if="id != 'new'">
    <q-bar class="bg-secondary text-white shadow-2">
      <div>{{label}}</div>
      <q-space />
      <q-btn dense flat icon="expand_less" v-if="isShowList && list.length > 0" @click="isShowList = false" />
      <q-btn dense flat icon="expand_more" v-if="!isShowList && list.length > 0" @click="isShowList = true" />
      <q-btn dense flat icon="add" @click="isShowAddDialog = true" v-if="!readonly"/>
    </q-bar>

    <q-list bordered separator v-if="isShowList" :dense="dense">
      <q-item v-for="item in filteredList" :key="item.id">
        <q-item-section avatar @click="$router.push(`${tableDependRoute}/${item[tableDependFldName]}`)">
          <q-avatar rounded>
            <comp-stat-img-src :src="avatar(item)"/>
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label>{{item.options.title[tableDependFldTitle]}}</q-item-label>
          <slot :item="item"></slot>
        </q-item-section>
        <q-item-section side v-if="!readonly">
          <div class="text-grey-8 q-gutter-xs">
            <q-btn v-if="flds" flat round icon="edit" size="sm" @click="showEditDialog(item)">
              <q-tooltip>Редактировать</q-tooltip>
            </q-btn>
            <q-btn flat round icon="delete" size="sm" @click="showDeleteDialog(item.id)">
              <q-tooltip>Удалить</q-tooltip>
            </q-btn>
          </div>
        </q-item-section>
      </q-item>
    </q-list>

    <!-- диалог добавления   -->
    <q-dialog v-model="isShowAddDialog" persistent>
      <q-card style="width: 700px; max-width: 80vw;">
        <q-card-section class="row items-center">
          <span class="q-ml-sm">Добавление</span>
        </q-card-section>

        <q-card-section>
          <comp-search-ref-in-list-widget label="Поиск" :pg-method="tableDependName + '_list'" @update="v => selectedForAdd = v" :ext="searchExt"/>
          <template v-if="flds">
            <div class="row q-col-gutter-md q-mt-sm" v-for="fldRow in flds">
              <comp-fld v-for="fld in fldRow" :key='fld.name'
                        :fld="localItem[fld.name]"
                        :type="fld.type"
                        @update="localItem[fld.name] = $event"
                        :label="fld.label"
                        :selectOptions="fld.selectOptions ? fld.selectOptions() : []"
                        :columnClass="fld.columnClass"
              />
            </div>
          </template>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Отмена" v-close-popup/>
          <q-btn flat label="Добавить" v-close-popup @click="add"/>
          <q-btn v-if="!hideCreateNew" flat label="Создать" v-close-popup @click="openNew" class="q-ml-md"/>
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="isShowEditDialog" persistent>
      <q-card style="width: 700px; max-width: 80vw;">
        <q-card-section>
          <template v-if="flds && selectedForEdit">
            <div class="row q-col-gutter-md q-mt-sm" v-for="fldRow in flds">
              <comp-fld v-for="fld in fldRow" :key='fld.name'
                        :fld="selectedForEdit[fld.name]"
                        :type="fld.type"
                        @update="selectedForEdit[fld.name] = $event"
                        :label="fld.label"
                        :selectOptions="fld.selectOptions ? fld.selectOptions() : []"
                        :columnClass="fld.columnClass"
              />
            </div>
          </template>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Отмена" v-close-popup/>
          <q-btn flat label="Ok" v-close-popup @click="save"/>
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- диалог подтверждения удаления   -->
    <q-dialog v-model="isShowDeleteDialog" persistent>
      <q-card>
        <q-card-section class="row items-center">
          <q-avatar rounded icon="warning" color="warning" text-color="white"/>
          <span class="q-ml-sm">Удалить?</span>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Отмена" v-close-popup/>
          <q-btn flat label="Удалить" v-close-popup @click="remove"/>
        </q-card-actions>
      </q-card>
    </q-dialog>

  </div>
</template>

<script>
  import _ from 'lodash'
  export default {
    props: ['id', 'tableIdName', 'tableIdFldName', 'tableDependName', 'tableDependFldName', 'tableDependRoute', 'linkTableName', 'label', 'avatarSrc', 'hideCreateNew', 'flds', 'readonly', 'searchExt', 'filterListFn', 'dense'],
    computed: {
      tableDependFldTitle() {
        return this.tableDependFldName.split('_id')[0] + '_title'
      },
      filteredList() {
        // возможность передать в компоненту функцию для фильтрации списка
        if (this.filterListFn) return this.list.filter(this.filterListFn)
        return this.list
      }
    },
    data() {
      return {
        isShowList: true,
        list: [],
        localItem: {},
        isShowAddDialog: false,
        isShowDeleteDialog: false,
        selectedForDeleteId: null,
        selectedForAdd: null,
        isShowEditDialog: false,
        selectedForEdit: null,
      }
    },
    methods: {
      avatar(item) {
        let fldName = this.tableDependFldName.split('_id')[0]
        return item.options.title[fldName + '_avatar'] || this.avatarSrc || 'https://www.svgrepo.com/show/95333/avatar.svg'
      },
      showDeleteDialog(id) {
        this.selectedForDeleteId = id
        this.isShowDeleteDialog = true
      },
      showEditDialog(v) {
        this.isShowEditDialog = true
        this.selectedForEdit = v
      },
      add() {
        if (!this.selectedForAdd) return
        let params = {id: -1}
        params[this.tableIdFldName] = +this.id
        params[this.tableDependFldName] = this.selectedForAdd.id
        this.$utils.postCallPgMethod({method: `${this.linkTableName}_update`, params: Object.assign(params, this.localItem)}).subscribe(res => {
          if (res.ok) {
            this.reload()
            this.$emit('reload')
          }
        })
      },
      remove() {
        this.$utils.postCallPgMethod({method: `${this.linkTableName}_update`, params: {id: this.selectedForDeleteId, deleted: true}}).subscribe(res => {
          if (res.ok) {
            this.reload()
            this.$emit('reload')
          }
        })
      },
      reload() {
        let params = {}
        params[this.tableIdFldName] = +this.id
        params.order_by = 'id'
        this.$utils.postCallPgMethod({method: `${this.linkTableName}_list`, params}).subscribe(res => {
          if (res.ok) {
            this.list = res.result
          }
        })
      },
      save() {
        this.$utils.postCallPgMethod({method: `${this.linkTableName}_update`, params: this.selectedForEdit}).subscribe(res => {
          if (res.ok) {
            this.isShowEditDialog = false
            this.selectedForEdit = null
            this.reload()
            this.$emit('reload')
          }
        })
      },
      openNew() {
        window.open(`${this.tableDependRoute}/new`, '_blank')
      }
    },
    mounted() {
      this.reload()
      if (this.flds) {
        _.flattenDeep(this.flds).map(v => {
          this.localItem[v.name] = v.type === 'checkbox' ? false : null
        })
      }
    }
  }
</script>
