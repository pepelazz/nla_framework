<template>
    <q-page padding>
        <comp-breadcrumb v-if="!isOpenInDialog && item"
                         :list="[{label: '[[index .Vue.I18n "listTitle"]]', to:`/[[.Vue.RouteName]]`,  docType: '[[.Name]]'},
                       {label: item ? (item.title ? item.title : 'Редактирование') : '',  docType: 'edit'}]"/>

        <div v-if="item" class="q-mt-sm">
            <!-- статус     -->
            <q-banner dense class="bg-grey-3 q-mb-md">
                <template v-slot:avatar>
                    <q-icon name="label" color="primary" />
                </template>
                <strong>{{item.title}}</strong>
                <q-chip>
                    <q-avatar color="secondary" icon="fas fa-star-half-alt"/>
                    {{$utils.i18n_[[.Name]]_state(item.state)}}
                    <q-tooltip>статус</q-tooltip>
                </q-chip>
            </q-banner>
            <!--  action кнопки   -->
            <div class="row q-gutter-sm q-mb-md" v-if="item.state != 'failed'">
                <div v-for="btn in actionBtnsByState[item.state]" :key="btn.compName" class="col-auto">
                    <component :is="btn.compName" :item="item" :currentUser="currentUser" @stateChanged="stateChanged"/>
                </div>
            </div>
            <div class="row q-col-gutter-sm">
                <!-- блок с параметрами сделки       -->
                <div class="col-xs-12 col-sm-8 col-md-8">

                    <!-- карточки с полями заполненными в стейтах     -->
                    <div v-for="(v, index) in stateList" :key="v.date" class="row q-col-gutter-md q-mb-sm">
                        <!-- в зависимости от стейта передаем в качестве item либо сам item, либо запись из истории item.options.states - потому что это не для редактирования, а для отображения                        -->
                        <component v-if="item" :is="'state_'+ v.state + '_card'" :state="item.state" :id="id" :item="item.state == v.state ? item : v" :date="v.date" :is_current_state="index === 0" :currentUser="currentUser" class="col-12"/>
                    </div>

                    <!--  кнопки   -->
                    <comp-item-btn-save @save="save" @cancel="$router.push(docUrl)"/>
                </div>
                <div class="col-xs-12 col-sm-4 com-md-4">
                    <comp-chat table_name="[[.Name]]" :table_id="id"/>
                </div>
<!--                    <q-tabs-->
<!--                            v-model="tab"-->
<!--                            dense-->
<!--                            class="bg-grey-3"-->
<!--                            align="justify"-->
<!--                            narrow-indicator-->
<!--                            inline-label-->
<!--                    >-->
<!--                        <q-tab name="chat" label="чат" icon="chat"/>-->
<!--                        <q-tab name="tasks" label="задачи" icon="list" />-->
<!--                    </q-tabs>-->
<!--                    <q-separator />-->
<!--                    <q-tab-panels v-model="tab" animated>-->
<!--                        <q-tab-panel name="chat">-->
<!--                            <comp-chat table_name="product_part_work" :table_id="id"/>-->
<!--                        </q-tab-panel>-->

<!--                        <q-tab-panel name="tasks">-->
<!--                            <task-list :id="id"/>-->
<!--                        </q-tab-panel>-->
<!--                    </q-tab-panels>-->
<!--                </div>-->
            </div>

        </div>
    </q-page>
</template>

<script>
    [[range .StateMachine.States]]
        [[- $state := . -]]
    import state_[[$state.Title]]_card from './comp/state_[[$state.Title]]_card'
    [[end]]
    [[- range .StateMachine.States]]
        [[- $state := . -]]
    [[- range .Actions]]
    import [[$state.Title]]_to_[[.To]]_btn from './comp/[[$state.Title]]_to_[[.To]]_btn'
    [[- end]]
    [[- end]]

    export default {
        props: ['id', 'isOpenInDialog', 'currentUser'],
        components: {[[range .StateMachine.States]] [[- $state := . -]] state_[[.Title]]_card, [[- range .Actions -]] [[$state.Title]]_to_[[.To]]_btn, [[end]] [[end]]},
        computed: {
            docUrl() {
                return `/[[.Vue.RouteName]]`
            }
        },
        data() {
            return {
                item: null,
                clientTitle: null,
                tab: 'chat',
                stateList: [],
                actionBtnsByState: {
                    [[range .StateMachine.States]]
                        [[- .Title]]:[
                            [[- $state := . -]]
                            [[- range .Actions]]
                            {compName: '[[$state.Title]]_to_[[.To]]_btn'},
                            [[- end]]
                        ],
                    [[end]]
                }
            }
        },
        methods: {
            resultModify(res) {
                [[.PrintVueItemResultModify]]
                return res
            },
            stateChanged() {
                this.$emit('reloadMsgList')
                this.reload()
            },
            save() {
                let itemForSave = this.item
                [[range .Flds ]]
                    [[ if eq .Vue.Type "select"]]itemForSave.[[.Name]] = this.item.[[.Name]] ? this.item.[[.Name]].value : null[[end]]
                [[- end]]
                this.$utils.postCallPgMethod({method: `[[.Name]]_update`, params: itemForSave}).subscribe(res => {
                    if (res.ok) {
                        this.item = this.resultModify(res.result)
                    }
                })
            },
            reload() {
                let cb = (v) => {
                    this.item = this.resultModify(v)
                    // заполняем последовательность стейтов
                    this.stateList = this.item.options.states
                }
                this.$utils.getDocItemById.call(this, {method: '[[.Name]]_get_by_id', cb})
            }
        },
        mounted() {
            this.reload()
        }
    }
</script>
