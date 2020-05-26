<template>
  <div v-if="id != 'new'">
    <q-bar class="bg-secondary text-white shadow-2">
      <div>{{label}}</div>
      <q-space />
      <q-btn dense flat icon="expand_less" v-if="isShowList && list.length > 0" @click="isShowList = false" />
      <q-btn dense flat icon="expand_more" v-if="!isShowList && list.length > 0" @click="isShowList = true" />
      <q-btn dense flat icon="add" @click="isShowAddDialog = true"/>
    </q-bar>

    <q-list bordered separator v-if="isShowList">
      <q-item v-for="item in list" :key="item.id">
        <q-item-section avatar @click="$router.push(`${tableDependRoute}/${item[tableDependName+'_id']}`)">
          <q-avatar rounded>
            <comp-stat-img-src :src="avatar(item)"/>
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label>{{item.options.title[tableDependFldTitle]}}</q-item-label>
        </q-item-section>
        <q-item-section side>
          <q-btn flat round icon="delete" size="sm" @click="showDeleteDialog(item.id)">
            <q-tooltip>Удалить</q-tooltip>
          </q-btn>
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
          <comp-search-ref-in-list-widget label="Поиск" :pg-method="tableDependName + '_list'" @update="v => selectedForAdd = v"/>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Отмена" v-close-popup/>
          <q-btn flat label="Добавить" v-close-popup @click="add"/>
          <q-btn flat label="Создать" v-close-popup @click="$router.push(`${tableDependRoute}/new`)" class="q-ml-md"/>
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
    export default {
        props: ['id', 'tableIdName', 'tableIdFldName', 'tableDependName', 'tableDependFldName', 'tableDependRoute', 'linkTableName', 'label', 'avatarSrc'],
        computed: {
            tableDependFldTitle() {
                return this.tableDependFldName.split('_id')[0] + '_title'
            },
            avatar() {
                return function(item) {
                    let fldName = this.tableDependFldName.split('_id')[0]
                    return item.options.title[fldName + '_avatar'] || this.avatarSrc || 'https://www.svgrepo.com/show/95333/avatar.svg'
                }
            },
        },
        data() {
            return {
                isShowList: true,
                list: [],
                isShowAddDialog: false,
                isShowDeleteDialog: false,
                selectedForDeleteId: null,
                selectedForAdd: null,
            }
        },
        methods: {
            showDeleteDialog(id) {
                this.selectedForDeleteId = id
                this.isShowDeleteDialog = true
            },
            add() {
                let params = {id: -1}
                params[this.tableIdFldName] = +this.id
                params[this.tableDependFldName] = this.selectedForAdd.id
                this.$utils.postCallPgMethod({method: `${this.linkTableName}_update`, params}).subscribe(res => {
                    if (res.ok) this.reload()
                })
            },
            remove() {
                this.$utils.postCallPgMethod({method: `${this.linkTableName}_update`, params: {id: this.selectedForDeleteId, deleted: true}}).subscribe(res => {
                    if (res.ok) {
                        this.reload()
                    }
                })
            },
            reload() {
                let params = {}
                params[`${this.tableIdName}_id`] = +this.id
                this.$utils.postCallPgMethod({method: `${this.linkTableName}_list`, params}).subscribe(res => {
                    if (res.ok) {
                        this.list = res.result
                    }
                })
            },
        },
        mounted() {
            this.reload()
        }
    }
</script>
