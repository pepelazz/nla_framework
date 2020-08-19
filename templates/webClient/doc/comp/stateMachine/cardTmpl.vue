<template>
    <div>
        <q-card flat bordered class="bg-grey-1">
            <q-item>
                <q-item-section avatar @click="isOpen = !isOpen">
                    <q-avatar rounded>
                        <img src="[[GetIconSrc]]">
                    </q-avatar>
                </q-item-section>
                <q-item-section>
                    <q-item-label>[[GetLabel]]</q-item-label>
                    <q-item-label caption>{{formatDateTime(date)}}</q-item-label>
                    <q-item-label caption v-if="item.deadline && !isReadonly"><q-badge outline color="primary" :label="formatDateTime(item.deadline)" /></q-item-label>
                </q-item-section>
                [[- if GetStateUpdateFldsGrid]]
                <q-item-section side>
                    <q-btn v-if="isOpen==false" icon="expand_more" round flat @click="isOpen=true"/>
                    <q-btn v-if="isOpen==true" icon="expand_less" round flat @click="isOpen=false"/>
                </q-item-section>
                [[- end]]
            </q-item>
            [[- if GetStateUpdateFldsGrid]]
            <q-separator/>
            <q-card-section v-if="isOpen">
                [[range GetStateUpdateFldsGrid]]
                [[- if .]]
                <div class="row q-col-gutter-md q-mb-sm">
                    [[range .]]
                    <div class='[[printf "%v" .Vue.ClassPrintOnlyCol]]'>
                        [[PrintVueFldTemplate .]]
                    </div>
                    [[end]]
                </div>
                [[- end -]]
                [[end]]
            </q-card-section>
            [[- end]]
        </q-card>
    </div>
</template>

<script>
    import isRole from '../../../mixins/isRole'
    export default {
        props: ['id', 'item', 'state', 'iconSrc', 'label', 'date', 'is_current_state', 'currentUser'],
        mixins: [isRole],
        components: {},
        computed: {
            isReadonly() {
                return !this.is_current_state
            },
            isManager() {
                return [this.item.manager_id, this.item.creator_id].includes(this.currentUser.id)
            },
            id: function () {
                return this.item?.id
            }
        },
        data() {
            return {
                isOpen: true,
            }
        },
        methods: {
            formatDateTime(d) {
                return this.$utils.formatPgDateTime(d)
            }
        },
        mounted() {
            [[.PrintVueItemStateMachineCardMounted]]
            this.isOpen = this.is_current_state
        }
    }
</script>
