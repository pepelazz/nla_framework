export default {
  data() {
    return {
      [[VarName]]: 0,
    }
  },
  mounted() {
    if (this.id !== 'new') {
      this.$utils.postCallPgMethod({method: '[[PgMethod]]', params: [[PgParams]]}).subscribe(res => {
        if (res.ok) {
          this.[[VarName]] = res.result.length
        }
      })
    }
  }
}
