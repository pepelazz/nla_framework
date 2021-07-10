import compDocList from './list/compDocList'
import compItemDropDownBtn from './list/itemDropDownBtn'
import compFld from './list/compFld'
import compFldUserSearch from './list/compFldUserSearch'
import compBreadcrumb from './list/compBreadcrumb'
import compDeleteBtnInList from './list/compDeleteBtnInList'
import compFldDate from './list/compFldDate'
import compFldDateTime from './list/compFldDateTime'
import compItemBtnSave from './list/compItemBtnSave'
import compFldUserMultipleSearch from './list/compFldUserMultipleSearch'
import compFldRefSearch from './list/compFldRefSearch'
import compFldSelectCity from './list/compFldSelectCity'
import compFldContact from './list/compFldContact'
import compFldFiles from './list/compFldFiles'
import compFldImg from './list/compFldImg'
import compFldImgList from './list/compFldImgList'
import compFldAddress from './list/compFldAddress'
import compDadataSuggestion from './list/compDadataSuggestion'
import compDadataAddress from './list/compDadataAddress'
import compDadataAddressDialog from './list/compDadataAddressDialog'
import compDadataCompany from './list/compDadataCompany'
import compSearchRefInListWidget from './list/compSearchRefInListWidget'
import compLinkListWidget from './list/compLinkListWidget'
import compFileUpload from './list/compFileUpload'
import compDialogConfirm from './list/compDialogConfirm'
import statImgSrc from './utils/statImgSrc'
import compDialogTaskAdd from './task/compDialogTaskAdd'
import compDialogTaskDone from './task/compDialogTaskDone'
import compChat from './chat/chat'

export default (Vue) => {
  Vue.component('comp-doc-list', compDocList)
  Vue.component('comp-item-dropdown-btn', compItemDropDownBtn)
  Vue.component('comp-fld', compFld)
  Vue.component('comp-fld-user-search', compFldUserSearch)
  Vue.component('comp-fld-user-multiple-search', compFldUserMultipleSearch)
  Vue.component('comp-breadcrumb', compBreadcrumb)
  Vue.component('comp-delete-btn-in-list', compDeleteBtnInList)
  Vue.component('comp-fld-date', compFldDate)
  Vue.component('comp-fld-date-time', compFldDateTime)
  Vue.component('comp-item-btn-save', compItemBtnSave)
  Vue.component('comp-fld-ref-search', compFldRefSearch)
  Vue.component('comp-fld-select-city', compFldSelectCity)
  Vue.component('comp-fld-contact', compFldContact)
  Vue.component('comp-fld-files', compFldFiles)
  Vue.component('comp-fld-img', compFldImg)
  Vue.component('comp-fld-address', compFldAddress)
  Vue.component('comp-fld-img-list', compFldImgList)
  Vue.component('comp-dadata-suggestion', compDadataSuggestion)
  Vue.component('comp-dadata-address', compDadataAddress)
  Vue.component('comp-dadata-address-dialog', compDadataAddressDialog)
  Vue.component('comp-dadata-company', compDadataCompany)
  Vue.component('comp-search-ref-in-list-widget', compSearchRefInListWidget)
  Vue.component('comp-link-list-widget', compLinkListWidget)
  Vue.component('comp-stat-img-src', statImgSrc)
  Vue.component('comp-dialog-confirm', compDialogConfirm)
  Vue.component('comp-dialog-task-add', compDialogTaskAdd)
  Vue.component('comp-dialog-task-done', compDialogTaskDone)
  Vue.component('comp-chat', compChat)
  Vue.component('comp-file-upload', compFileUpload)
}
