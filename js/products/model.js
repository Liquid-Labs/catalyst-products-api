import {
  Address,
  arrayType,
  CommonResourceConf,
  Model,
  entityPropsModel
} from '@liquid-labs/catalyst-core-api'

const productPropsModel = [
  'displayName',
  'summary',
  'supportPhone',
  'supportEmail',
  'homepage',
  'logoURL',
  'repoURL',
  'issuesURL']
  .map((propName) => ({ propName : propName, writable : true }))
productPropsModel.push(...entityPropsModel)
productPropsModel.push({
  propName  : 'addresses',
  model     : Address,
  valueType : arrayType,
  writable  : true})
productPropsModel.push({
  propName            : 'changeDesc',
  unsetForNew         : true,
  writable            : true,
  optionalForComplete : true
})

const Product = class extends Model {
  get resourceName() { return 'products' }
}
Model.finalizeConstructor(Product, productPropsModel)

const productResourceConf = new CommonResourceConf('product', {
  model       : Product,
  sortOptions : [
    { label : 'Dispaly name (asc)',
      value : 'displayName-asc',
      func  : (a, b) => a.displayName.localeCompare(b.displayName) },
    { label : 'Display name (desc)',
      value : 'displayName-desc',
      func  : (a, b) => -a.displayName.localeCompare(b.displayName) }
  ],
  sortDefault : 'displayName-asc'
})

export { Product, productPropsModel, productResourceConf }
