<template>
    <div>

        <div v-if="item" class="q-mt-sm">
            <!--  поля формы    -->
            [[ define "vueItemRow1" ]]
            <div class="[[.Class]]">
                [[- if gt (len .Grid) 0]]
                [[- range .Grid -]]
                [[ template "vueItemRow1" . ]]
                [[- end]]
                [[- else]]
                [[PrintVueFldTemplate .Fld]]
                [[- end]]
            </div>
            [[- end -]]
            [[range .Vue.Grid]]
            [[- template "vueItemRow1" .]]
            [[end]]

            [[- if .IsRecursion -]]
            <div class="row q-col-gutter-md q-mb-sm q-mt-sm">
                <div class="col-8" v-if="id !== 'new'">
                    <comp-recursive-child-list :id='id' @update='save'/>
                </div>
            </div>
            [[end]]

            <!--  кнопки   -->
            <comp-item-btn-save v-if="!isOpenInDialog" @save="save" @cancel="$router.push(docUrl)"/>
            <!--  при открытии в диалоге кнопку Отмена не показываем   -->
            <q-btn v-else color="secondary" label="сохранить" class="q-mr-sm" @click="save"/>

            [[range .Vue.Hooks.ItemHtml]]
            [[.]]
            [[- end]]

        </div>
    </div>
</template>

<script>
    [[ .PrintVueImport "docItem" ]]
    export default {
        props: ['id', 'isOpenInDialog' [[- if .IsRecursion -]], 'parent_id'[[- end -]]],
        components: {[[- .PrintComponents "docItem" -]]},
    mixins: [ [[- .Vue.PrintMixins "docItem" -]] ],
    computed: {
        docUrl: function() {
            return [[if not .IsRecursion -]]'/[[.Vue.RouteName]]'[[else -]] this.parent_id ? `/[[.Vue.RouteName]]/${this.parent_id}` : '/[[.Vue.RouteName]]' [[- end]]
        },
    },
    data() {
        return {
            item: null,
            flds: [
                [[- range .Flds]]
                    [[- if .Name]]
                {name: '[[.Name]]', label: '[[.Vue.NameRu]]'[[if .Vue.IsRequired -]],  required: true[[- end]]},
                [[- end]]
        [[- end]]
    ],
        optionsFlds: [ [[- .PrintVueItemOptionsFld -]] ],
    }
    },
    methods: {
        [[ .PrintVueMethods "docItem" ]]
        resultModify(res) {
            [[.PrintVueItemResultModify]]
            return res
        },
        save() {
            this.$utils.saveItem.call(this, {
                method: '[[.PgName]]_update',
                itemForSaveMod: {[[.PrintVueItemForSave]]},
            resultModify: this.resultModify,
        })
        },
    },
    mounted() {
        let cb = (v) => {
            this.item = this.resultModify(v)
        }
        this.$utils.getDocItemById.call(this, {method: '[[.PgName]]_get_by_id', cb})
    }
    }
</script>
