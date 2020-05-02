import compDocList from './list/compDocList'
import compItemDropDownBtn from './list/itemDropDownBtn'
import compFld from './list/compFld'
import compFldUserSearch from './list/compFldUserSearch'
import compBreadcrumb from './list/compBreadcrumb'
import compFldDate from './list/compFldDate'
import compFldDateTime from './list/compFldDateTime'
import compItemBtnSave from './list/compItemBtnSave'
import compFldUserMultipleSearch from './list/compFldUserMultipleSearch'
import compFldRefSearch from './list/compFldRefSearch'
import compFldSelectCity from './list/compFldSelectCity'
import compFldContact from './list/compFldContact'
import compDadataSuggestion from './list/compDadataSuggestion'
import compDadataAddress from './list/compDadataAddress'
import compSearchRefInListWidget from './list/compSearchRefInListWidget'
import compLinkListWidget from './list/compLinkListWidget'

export default (Vue) => {
  Vue.component('comp-doc-list', compDocList)
  Vue.component('comp-item-dropdown-btn', compItemDropDownBtn)
  Vue.component('comp-fld', compFld)
  Vue.component('comp-fld-user-search', compFldUserSearch)
  Vue.component('comp-fld-user-multiple-search', compFldUserMultipleSearch)
  Vue.component('comp-breadcrumb', compBreadcrumb)
  Vue.component('comp-fld-date', compFldDate)
  Vue.component('comp-fld-date-time', compFldDateTime)
  Vue.component('comp-item-btn-save', compItemBtnSave)
  Vue.component('comp-fld-ref-search', compFldRefSearch)
  Vue.component('comp-fld-select-city', compFldSelectCity)
  Vue.component('comp-fld-contact', compFldContact)
  Vue.component('comp-dadata-suggestion', compDadataSuggestion)
  Vue.component('comp-dadata-address', compDadataAddress)
  Vue.component('comp-search-ref-in-list-widget', compSearchRefInListWidget)
  Vue.component('comp-link-list-widget', compLinkListWidget)
}
