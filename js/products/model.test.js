/* globals beforeAll describe expect test */
import { resourcesSettings, verifyCatalystSetup } from '@liquid-labs/catalyst-core-api'
import { Product, productResourceConf } from './model'

const productFooModel = {
  pubId       : '630AC9ED-3531-41E3-BD87-E26ADA74ECBC',
  lastUpdated : null,
  displayName : 'foo',
  summary     : null,
  supportPhone       : null,
  supportEmail       : null,
  homepage    : null,
  logoURL     : null,
  repoURL : null,
  issuesURL : null,
  ontology : 'TANGIBLE GOOD',
  addresses   : undefined
}

const productBarModel = {
  pubId       : '23DB5195-67FF-4709-9033-7F9F5C5A6C6F',
  lastUpdated : null,
  displayName : 'bar',
  summary     : null,
  supportPhone       : null,
  supportEmail       : null,
  homepage    : null,
  logoURL     : null,
  repoURL : null,
  issuesURL : null,
  ontology : 'DIGITAL GOOD',
  addresses   : []
}

describe('Product', () => {
  beforeAll(() => {
    const resourceList = [ productResourceConf ]
    resourcesSettings.setResources(resourceList)
    verifyCatalystSetup()
  })

  test("should identify self as a 'products' resource", () => {
    const product = new Product(productFooModel)
    expect(product.resourceName).toBe('products')
  })

  test("should be incomplete if address is 'null'", () => {
    const product = new Product(productFooModel)
    expect(product.isComplete()).toBe(false)
    expect(product.getMissing()).toHaveLength(1)
    expect(product.getMissing()[0]).toBe('addresses')
  })

  test("should provide ascending and descending display name sort options", () => {
    const productFoo = new Product(productFooModel)
    const productBar = new Product(productBarModel)

    const products = [ productFoo, productBar ]
    expect(typeof resourcesSettings.getResourcesMap()['products'].sortMap['displayName-asc'])
      .toBe('function')
    products.sort(resourcesSettings.getResourcesMap()['products'].sortMap['displayName-asc'])
    expect(products[0]).toBe(productBar)
    expect(products[1]).toBe(productFoo)

    expect(typeof resourcesSettings.getResourcesMap()['products'].sortMap['displayName-desc'])
      .toBe('function')
    products.sort(resourcesSettings.getResourcesMap()['products'].sortMap['displayName-desc'])
    expect(products[0]).toBe(productFoo)
    expect(products[1]).toBe(productBar)
    // and verify that we test all the options
    expect(resourcesSettings.getResourcesMap()['products'].sortOptions).toHaveLength(2)
  })

  test("should define default sort options", () => {
    expect(resourcesSettings.getResourcesMap()['products'].sortDefault).toBe('displayName-asc')
  })
})
