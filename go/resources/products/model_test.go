package products_test

import (
  "encoding/json"
  "reflect"
  "strings"
  "testing"

  . "github.com/Liquid-Labs/catalyst-products-api/go/resources/products"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/entities"
  "github.com/Liquid-Labs/go-nullable-mysql/nulls"
  "github.com/stretchr/testify/assert"
)

var widgetProduct = &Product{
  entities.Entity{
    nulls.NewInt64(1),
    nulls.NewString(`a`),
    nulls.NewInt64(2),
  },
  nulls.NewString(`4C2B3954-8D7F-48BA-B720-3B0F15F91BA9`),
  nulls.NewString(`Widget`),
  nulls.NewString(`A better dodad.`),
  nulls.NewString(`foo@test.com`),
  nulls.NewString(`555-555-9999`),
  nulls.NewString(`https://foo.com/products/widget`),
  nulls.NewString(`http://foo.com/assets/widget_logo.svg`),
  nulls.NewString(`https://foo.com/products/widget/repo`),
  nulls.NewString(`https://foo.com/products/widget/issues`),
  nulls.NewString(`TANGIBLE GOOD`),
}

func TestProductClone(t *testing.T) {
  clone := widgetProduct.Clone()
  assert.Equal(t, widgetProduct, clone, `Original does not match clone.`)
  clone.Id = nulls.NewInt64(3)
  clone.PubId = nulls.NewString(`b`)
  clone.LastUpdated = nulls.NewInt64(4)
  clone.SetLegalOwnerPubID(`D6F0B9B4-078F-49E8-9905-1E399B6A9AEF`)
  clone.SetDisplayName(`different name`)
  clone.SetSummary(`A new summary.`)
  clone.SetSupportEmail(`blah@test.com`)
  clone.SetSupportPhone(`555-555-9997`)
  clone.SetHomepage(`https://bar.com`)
  clone.SetLogoURL(`http://bar.com/image`)
  clone.SetRepoURL(`https://bar.com/widget_repo`)
  clone.SetIssuesURL(`https://bar.com/issues`)
  clone.SetOntology(`DIGITAL GOOD`)

  oReflection := reflect.ValueOf(widgetProduct).Elem()
  cReflection := reflect.ValueOf(clone).Elem()
  for i := 0; i < oReflection.NumField(); i++ {
    assert.NotEqualf(
      t,
      oReflection.Field(i).Interface(),
      cReflection.Field(i).Interface(),
      `Fields '%s' unexpectedly match with: %s`,
      oReflection.Type().Field(i).Name, oReflection.Field(i).Interface(),
    )
  }
}

const jsonDisplayName = "A Product"
const jsonEmail = "johndoe@test.com"
const jsonPhone = "555-555-0000"
const jsonHomepage = `https://foo.com`

var johnDoeJson string = `
  {
    "displayName": "` + jsonDisplayName + `",
    "supportEmail": "` + jsonEmail + `",
    "supportPhone": "` + jsonPhone + `",
    "homepage": "` + jsonHomepage + `"
  }`

var decoder *json.Decoder = json.NewDecoder(strings.NewReader(johnDoeJson))
var someProduct = &Product{}
var decodeErr = decoder.Decode(someProduct)

func TestProductsDecode(t *testing.T) {
  assert.NoError(t, decodeErr, "Unexpected error decoding product JSON.")
  assert.Equal(t, jsonDisplayName, someProduct.DisplayName.String, "Unexpected display name.")
  assert.Equal(t, jsonEmail, someProduct.SupportEmail.String, "Unexpected email.")
  assert.Equal(t, jsonPhone, someProduct.SupportPhone.String, "Unexpected phone.")
  assert.Equal(t, jsonHomepage, someProduct.Homepage.String, "Unexpected homepage value.")
}

func TestProductFormatter(t *testing.T) {
  testP := &Product{ SupportPhone: nulls.NewString(`5555555555`) }
  testP.FormatOut()
  assert.Equal(t, `555-555-5555`, testP.SupportPhone.String)
}
