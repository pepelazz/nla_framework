<template>
  <q-drawer :modelValue="leftSide" side="left" bordered show-if-above
            content-class="bg-grey-2"
            :width="240" @hide="$emit('hide')">
<!--    <q-scroll-area class="fit">-->
      <q-list padding>
        <div v-for="link in menuLinks" :key="link.text">
            <span v-if="isRole(link.roles, link.conditionFunc)">
              <q-item v-if="!link.isFolder" v-ripple clickable :to="link.url" exact>
                <q-item-section avatar>
                  <q-avatar rounded>
                    <img :src="link.icon">
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label v-html='link.text.includes("i18n_") ? $t(link.text.replace("i18n_", "")) : link.text'/>
                </q-item-section>
              </q-item>
              <q-expansion-item v-else header-class="text-grey" closed>
                 <template v-slot:header>
                  <q-item-section avatar>
                    <q-avatar rounded>
                      <img :src="link.icon">
                    </q-avatar>
                  </q-item-section>
                    <q-item-section>
                      <q-item-label v-html='link.text.includes("i18n_") ? $t(link.text.replace("i18n_", "")) : link.text'/>
                    </q-item-section>
                </template>
                <span v-for="link1 in link.linkList" :key="link1.text" >
                  <q-item v-if="isRole(link1.roles, link1.conditionFunc)" :inset-level="1" v-ripple clickable :to="link1.url"
                          exact>
                  <q-item-section>
                    <q-item-label v-html='link1.text.includes("i18n_") ? $t(link1.text.replace("i18n_", "")) : link1.text'/>
                  </q-item-section>
                </q-item>
                </span>
              </q-expansion-item>
            </span>
        </div>
      </q-list>
<!--    </q-scroll-area>-->
  </q-drawer>
</template>

<script>
    export default {
        props: ['leftSide', 'currentUser'],
        computed: {
          isRole() {
            return (roles, conditionFunc) => {
                  if (conditionFunc) {
                    let isAccess = conditionFunc(this.currentUser)
                    console.log('isAccess:', isAccess)
                    if (!isAccess) return false
                  }
                  if (!roles || roles.length === 0) return true
                  if (!this.currentUser.role) return false
                  let isAccess = false
                  this.currentUser.role.map(r => {
                      if (roles.includes(r)) isAccess = true
                  })
                  return isAccess
              }
          }
        },
        data() {
            return {
                //
                menuLinks: [
                    // for codeGenerate ##sidemenu_slot1
                ],
            }
        },
    }
</script>
