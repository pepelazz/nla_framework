<template>
  <div class="row">
    <div :class="colClassComputed">
      <q-infinite-scroll @load="itemListLoad" :offset="250" ref="infiniteScroll">
        <!-- список         -->
        <q-list bordered class="rounded-borders" separator>
          <q-item>
            <q-item-section>
              <q-item-label caption>{{computedListTitle}}</q-item-label>
            </q-item-section>
            <q-item-section top side>
              <div class="text-grey-8 q-gutter-xs">
                <!-- кнопка создания нового документа   -->
                <template v-if="!readonly && (newDocUrl || newDocEventOnly)">
                  <q-btn v-if="!newDocEventOnly" size="12px" flat dense round icon="add" @click="openNewDoc"/>
                  <q-btn v-else size="12px" flat dense round icon="add" @click="$emit('clickAddBtn')"/>
                </template>
                <!-- дополнительные кнопки   -->
                <slot name="addBtnsSlot"/>
                <!-- кнопка поиска  -->
                <q-btn size="12px" v-if="searchFldName" flat dense round icon="search" @click="toggleSearchFld"/>
                <!--  СОРТИРОВКА ПО ВОЗРАСТАНИЮ  -->
                <q-btn v-if="listSortData" size="12px" flat dense round icon="expand_less">
                  <q-menu auto-close>
                    <q-list dense style="min-width: 100px">
                      <q-item clickable v-for="item in listSortData" :key="item.value">
                        <q-item-section @click="changeItemList({order_by: `${item.value}`})">{{item.title}}
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
                <!--  СОРТИРОВКА ПО УБЫВАНИЮ  -->
                <q-btn v-if="listSortData" size="12px" flat dense round icon="expand_more">
                  <q-menu auto-close>
                    <q-list dense style="min-width: 100px">
                      <q-item clickable v-for="item in listSortData" :key="item.value">
                        <q-item-section @click="changeItemList({order_by: `${item.value} desc`})">{{item.title}}
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
                <!--  ФИЛЬТР  -->
                <q-btn v-if="listFilterData" size="12px" flat dense round icon="filter_list">
                  <q-menu auto-close>
                    <q-list dense style="min-width: 100px">
                      <q-item clickable v-for="item in listFilterData" :key="item.title">
                        <q-item-section @click="changeItemList(item.value)">{{item.title}}
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
              </div>
            </q-item-section>
          </q-item>

          <!--  поле поиска  -->
          <q-item v-if="isShowSearchfld && searchFldName">
            <q-item-section top>
              <div class="text-grey-8 q-gutter-xs">
                <q-input ref="searchInput" dense filled v-model="searchTxt" class="q-ml-md">
                  <template v-slot:append>
                    <q-icon v-if="searchTxt === ''" name="search"/>
                    <q-icon v-else name="clear" class="cursor-pointer" @click="searchTxt = ''"/>
                  </template>
                </q-input>
              </div>
            </q-item-section>
          </q-item>

          <!-- дополнительные поля (фильтр по датам и пр)   -->
          <slot name="addFilterSlot"/>

          <q-item v-for="item in itemList" :key="item.id">
            <!--  слот для рендеринга элемента списка -->
            <slot name="listItem" :item="item"></slot>
          </q-item>

        </q-list>
        <!-- спинер загрузки   -->
        <template v-slot:loading>
          <div class="row justify-center q-my-md">
            <q-spinner-dots color="primary" size="40px"/>
          </div>
        </template>
      </q-infinite-scroll>
    </div>
    <!-- кнопка создания нового документа   -->
    <q-page-sticky v-if="!readonly && (newDocUrl || newDocEventOnly)" position="bottom-right" :offset="[18, 18]">
      <q-btn v-if="!newDocEventOnly" fab icon="add" color="accent" @click="openNewDoc"/>
      <q-btn v-else fab icon="add" color="accent" @click="$emit('clickAddBtn')"/>
    </q-page-sticky>
  </div>
</template>

<script>
  import {debounce} from 'quasar'
  import _ from 'lodash'

  export default {
    props: ['listTitle','listDeletedTitle', 'pgMethod', 'listSortData', 'listFilterData', 'searchFldName', 'newDocEventOnly', 'newDocUrl', 'isOpenNewInTab', 'urlQueryParams', 'ext', 'readonly', 'colClass', 'startFilter'],
    computed: {
      computedListTitle() {
        return !this.listParams.deleted ? this.listTitle : this.listDeletedTitle
      },
      colClassComputed() {
        return this.colClass || 'col-xs-12 col-sm-12 col-md-6 q-gutter-md q-pt-md'
      }
    },
    data() {
      return {
        searchTxt: '',
        isShowSearchfld: false,
        itemList: [],
        listParams: {page: 0, per_page: 10, deleted: false},
        isUrlQueryProcessed: false, // флаг для обработки query из url при первоначальной загрузке
      }
    },
    methods: {
      toggleSearchFld() {
        this.isShowSearchfld = !this.isShowSearchfld
        if (!this.isShowSearchfld) this.listParams[this.searchFldName] = ''
        if (this.isShowSearchfld) {
          this.$nextTick(() => {
            this.$refs.searchInput.focus()
          })
        }
      },
      itemListLoad(index, done) {
        // при первоначальной загрузке обрабатываем query из url
        if (!this.isUrlQueryProcessed) {
          this.isUrlQueryProcessed = true
          if (this.urlQueryParams) {
            const urlParams = new URLSearchParams(window.location.search)
            this.urlQueryParams.map(v => {
              // объекты обрабатываем отдельно, например даты. Можно для них здесь написать отдельную обработку
              if (typeof v === 'string' && urlParams.has(v)) {
                this.listParams[v] = urlParams.get(v)
              }
            })
          }
        }
        this.loadList({list: this.itemList, params: this.listParams, done})
      },
      loadList({list = [], params = {}, done}) {
        this.listParams.page++
        // обновляем параметры в query параметрах списка
        this.$utils.updateUrlQuery(_.omit(params, ['per_page', 'page']))
        this.$utils.postCallPgMethod({method: this.pgMethod, params: Object.assign(params, this.ext ? this.ext : {})}).subscribe(res => {
          if (res.ok) {
            if (res.result && res.result.length > 0) {
              res.result.map(v => list.push(v))
              this.$emit('updateCount', list.length)
              if (done) done()
            } else {
              if (done) done(true)
            }
          } else {
            if (done) done(true)
          }
        })
      },
      reloadList() {
        this.itemList = []
        this.listParams.page = 0
        this.$refs.infiniteScroll.resume()
        this.loadList({list: this.itemList, params: this.listParams})
        this.$forceUpdate()
      },
      reloadListDebounce() {
        this.reloadList()
      },
      changeItemList(params) {
        this.listParams = Object.assign(this.listParams, params)
        this.reloadList()
      },
      openNewDoc() {
        // либо открываем в новом табе, либо в этом же
        this.isOpenNewInTab ? window.open('/' + this.newDocUrl, '_blank') : this.$router.push(this.newDocUrl)
      }
    },
    watch: {
      searchTxt(v) {
        if (this.searchFldName && (v.length === 0 || v.length > 3)) {
          this.listParams[this.searchFldName] = v
          this.reloadListDebounce()
        }
      }
    },
    mounted() {
      // если указан начальный фильтр, то применяем его при первоначальной загрузке
      if (this.startFilter) {
        this.listParams = Object.assign(this.listParams, this.startFilter)
      }
      this.loadList({list: this.itemList, params: this.listParams})
      this.reloadListDebounce = debounce(this.reloadListDebounce, 300)
    }
  }
</script>
