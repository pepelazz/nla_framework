export default [
  {label: 'сотрудник', value: 'student'},
  {label: 'админ', value: 'admin'},
  [[- range .Roles]]
  {label: '[[.NameRu]]', value: '[[.Name]]'},
  [[- end]]
]
