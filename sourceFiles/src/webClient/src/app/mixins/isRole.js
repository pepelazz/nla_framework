export default {
  computed: {
    isRole() {
      return (...roles) => {
        // из [[]] -> []
        roles = this.$utils._.flatMap(roles)
        if (!roles || roles.length === 0) return true
        if (!this.currentUser.role) return false
        let isAccess = false
        this.currentUser.role.map(r => {
          if (roles.includes(r)) isAccess = true
        })
        return isAccess
      }
    }
  }
}
