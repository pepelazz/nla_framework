import {Subject} from 'rxjs'
import {takeUntil} from 'rxjs/operators'

export default {
  computed: {
    isLoggedIn() {
      return this.currentUser
    },
    isWorking() {
      return this.currentUser && this.currentUser.options && this.currentUser.options.state === 'working'
    },
    isWaitingAuth() {
      return this.currentUser && this.currentUser.options && this.currentUser.options.state === 'waiting_auth'
    },
    isFired() {
      return this.currentUser && this.currentUser.options && this.currentUser.options.state === 'fired'
    },
    isAdmin() {
      return this.currentUser?.role?.includes('admin')
    }
  },
  data() {
    return {
      currentUser: {},
      isInLogingProcess: false,
      currentUserdestroy$: new Subject(),
    }
  },
  methods: {
    login() {
      this.$currentUser.login()
    },
    logout() {
      this.$currentUser.logout()
    },
  },
  mounted() {
    this.$currentUser.getUser$().pipe(takeUntil(this.currentUserdestroy$)).subscribe(v => {
      this.currentUser = v
    })
    this.$currentUser.getIsInLogingProcess().pipe(takeUntil(this.currentUserdestroy$)).subscribe(v => {
      this.isInLogingProcess = v
    })
  },
  destroyed() {
    this.currentUserdestroy$.next(true)
    this.currentUserdestroy$.complete()
  }
}
