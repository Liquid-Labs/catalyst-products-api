package products_test

import (
  "context"
  "os"
  "testing"

  // the package we're testing
  . "github.com/Liquid-Labs/catalyst-products-api/go/resources/products"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/entities"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/locations"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/users"
  "github.com/Liquid-Labs/go-api/sqldb"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
)

func TestProductsDBIntegration(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  }

  if widgetProduct == nil {
    t.Error(`Product struct not define; can't continue. This probbaly indicates a setup failure in 'model_test.go'.`)
  } else {
    if t.Run(`ProductsDBSetup`, testProductDBSetup) {
      if sqldb.DB == nil { // test was skipped, but we still need to setup
        setupDB()
      }
      t.Run(`ProductGet`, testProductGet)
      t.Run(`ProductCreate`, testProductCreate)
      t.Run(`ProductUpdate`, testProductUpdate)
      t.Run(`ProductGetInTxn`, testProductGetInTxn)
      t.Run(`ProductCreateInTxn`, testProductCreateInTxn)
      t.Run(`ProductUpdateInTxn`, testProductUpdateInTxn)
    }
  }
}

const someProductID=`D929BEE3-8034-40A9-B33E-E1A28507EE68`

func setupDB() {
  sqldb.RegisterSetup(entities.SetupDB, locations.SetupDB, users.SetupDB, /*products.*/SetupDB)
  sqldb.InitDB() // panics if unable to initialize
}

func testProductDBSetup(t *testing.T) {
  setupDB()
}

func testProductGet(t *testing.T) {
  product, err := GetProduct(someProductID, context.Background())
  require.NoError(t, err, `Unexpected error getting Product.`)
  require.NotNil(t, product, `Unexpected nil Product on create (with no error); check ID.`)
  assert.Equal(t, `Bauble`, product.DisplayName.String, `Unexpected display name.`)
  assert.Equal(t, `A thing for your wall.`, product.Summary.String, `Unexpected summary.`)
  assert.Equal(t, `bauble@foo.com`, product.SupportEmail.String, `Unexpected email.`)
  assert.False(t, product.SupportPhone.Valid, `Unexpected phone.`)
  assert.Equal(t, `https://foo.com/proudcts/bauble`, product.Homepage.String, `Unexpected homepage value.`)
  assert.NotEmpty(t, product.Id, `Unexpected empty ID.`)
  assert.Equal(t, someProductID, product.PubId.String, `Unexpected public id.`)
}

func testProductCreate(t *testing.T) {
  product, err := CreateProduct(widgetProduct, context.Background())
  require.NoError(t, err, `Unexpected error creating Product.`)
  require.NotNil(t, product, `Unexpected nil Product on create (with no error).`)
  assert.Equal(t, widgetProduct.DisplayName, product.DisplayName, `Unexpected display name.`)
  assert.Equal(t, widgetProduct.SupportEmail, product.SupportEmail, `Unexpected email.`)
  assert.Equal(t, widgetProduct.SupportPhone, product.SupportPhone, `Unexpected phone.`)
  assert.Equal(t, widgetProduct.Homepage, product.Homepage, `Unexpected homepage value.`)
  assert.NotEmpty(t, product.Id, `Unexpected empty ID.`)
  assert.NotEmpty(t, product.PubId, `Unexpected empty public id.`)
}

func testProductUpdate(t *testing.T) {
  someOtherProduct, err := GetProduct(someProductID, context.Background())
  require.NoError(t, err, `Unexpected error getting Product.`)
  require.NotNil(t, someOtherProduct, `Unexpected nil value retrieving product (with no error); check ID.`)
  someOtherProduct.SetDisplayName(`Bauble 2.0`)
  someOtherProduct.SetSupportEmail(`janepdoe@test.com`)
  someOtherProduct.SetSupportPhone(`555-555-0001`)
  someOtherProduct.SetHomepage(`https://producthome.com`)
  product, err := UpdateProduct(someOtherProduct, context.Background())
  require.NoError(t, err, `Unexpected error updating Product.`)
  require.NotNil(t, product, `Unexpected nil Product on create (with no error).`)
  assert.Equal(t, someOtherProduct.DisplayName, product.DisplayName, `Unexpected display name.`)
  assert.Equal(t, someOtherProduct.SupportEmail, product.SupportEmail, `Unexpected email.`)
  assert.Equal(t, someOtherProduct.SupportPhone, product.SupportPhone, `Unexpected phone.`)
  assert.Equal(t, someOtherProduct.Homepage, product.Homepage, `Unexpected active value.`)
  assert.NotEmpty(t, product.Id, `Unexpected empty ID.`)
  assert.NotEmpty(t, product.PubId, `Unexpected empty public id.`)
}

func testProductGetInTxn(t *testing.T) {
  someOtherProduct, restErr := GetProduct(someProductID, context.Background())
  assert.NoError(t, restErr, `Unexpected error getting product.`)
  txn, _ := sqldb.DB.Begin()
  orig := someOtherProduct.Clone()
  // if we get in a txn, we should see the changes
  someOtherProduct.SetSupportPhone(`555-555-0003`)
  product, restErr := UpdateProductInTxn(someOtherProduct, context.Background(), txn)
  someOtherTxn, restErr := GetProductInTxn(someProductID, context.Background(), txn)
  assert.Equal(t, *product, *someOtherTxn, `Update-Product and Get-Product do not match.`)
  assert.Equal(t, someOtherProduct.SupportPhone, someOtherTxn.SupportPhone, `Did not see change while getting in txn.`)
  assert.NotEqual(t, someOtherProduct.SupportPhone, orig.SupportPhone, `Phone number not changed.`)
  someOtherNoTxn, restErr := GetProduct(someProductID, context.Background())
  assert.Equal(t, orig.SupportPhone, someOtherNoTxn.SupportPhone, `Non-txn product reflects changes.`)
  assert.NoError(t, txn.Commit(), `Error attempting commit.`)
  someOtherFinish, _ := GetProduct(someProductID, context.Background())
  assert.Equal(t, *someOtherTxn, *someOtherFinish, `Post-commit Products didn't match.`)
}

func testProductCreateInTxn(t *testing.T) {
  yetAnotherProduct := widgetProduct.Clone()
  yetAnotherProduct.SetDisplayName(`Bauble 3.0`)
  txn, _ := sqldb.DB.Begin()
  txnProduct, restErr := CreateProductInTxn(yetAnotherProduct, context.Background(), txn)
  assert.NoError(t, restErr, `Unexpected error creating product in txn.`)
  noProduct, restErr := GetProduct(txnProduct.PubId.String, context.Background())
  assert.Nil(t, noProduct, `Unexpected retrieval of product outside of txn.`)
  assert.Error(t, restErr, `Unexpected non-error while retrieving product outside of txn.`)
  assert.NoError(t, txn.Commit(), `Error attempting commit.`)
  yetAnotherFinish, _ := GetProduct(txnProduct.PubId.String, context.Background())
  assert.Equal(t, *txnProduct, *yetAnotherFinish, `Post-commit Products didn't match.`)
}

func testProductUpdateInTxn(t *testing.T) {
  /*txn, err := sqldb.DB.Begin()
  assert.NoError(t, err, `Unexpected error opening transaction.`)*/
}
