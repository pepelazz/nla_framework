
import {ref, onMounted, inject, getCurrentInstance, onBeforeUnmount} from 'vue'
import {takeUntil} from 'rxjs/operators'
import {Subject} from 'rxjs'

const getCurrentUser = (currentUser) => {
  onMounted(()=>{
    let currentUserdestroy$ = new Subject()
    const app = getCurrentInstance()
    app.appContext.config.globalProperties.$currentUser.getUser$().pipe(takeUntil(currentUserdestroy$)).subscribe(v => {
      currentUser.value = v
    })
    onBeforeUnmount(()=>{
      currentUserdestroy$.next(true)
      currentUserdestroy$.complete()
    })
  })
}

export default getCurrentUser
