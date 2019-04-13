package products

import (
  "context"
  "database/sql"
  "fmt"
  "log"

  "github.com/Liquid-Labs/go-api/sqldb"
  "github.com/Liquid-Labs/go-nullable-mysql/nulls"
  "github.com/Liquid-Labs/go-rest/rest"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/entities"
)

var ProductsSorts = map[string]string{
  "": `p.display_name ASC `,
  `name-asc`: `p.display_name ASC `,
  `name-desc`: `p.display_name DESC `,
}

func ScanProduct(row *sql.Rows) (*Product, error) {
	var p Product

	if err := row.Scan(&p.Id, &p.PubId, &p.LastUpdated, &p.LegalOwnerPubID, &p.DisplayName, &p.Summary, &p.SupportPhone, &p.SupportEmail, &p.Homepage, &p.LogoURL, &p.RepoURL, &p.IssuesURL, &p.Ontology); err != nil {
		return nil, err
	}

	return &p, nil
}

// implement rest.ResultBuilder
func BuildProductResults(rows *sql.Rows) (interface{}, error) {
  results := make([]*Product, 0)
  for rows.Next() {
    product, err := ScanProduct(rows)
    if err != nil {
      return nil, err
    }

    results = append(results, product)
  }

  return results, nil
}

// Implements rest.GeneralSearchWhereBit
func ProductsGeneralWhereGenerator(term string, params []interface{}) (string, []interface{}, error) {
  likeTerm := `%`+term+`%`
  var whereBit string = "AND (p.display_name LIKE ? OR p.summary LIKE ?) "
  params = append(params, likeTerm, likeTerm)

  return whereBit, params, nil
}

const CommonProductFields = `e.id, e.pub_id, e.last_updated, lo.pub_id, p.display_name, p.summary, p.support_phone, p.support_email, p.homepage, p.logo_url, p.repo_url, p.issues_url, p.ontology `
const CommonProductsFrom = `FROM products p JOIN entities e ON p.id=e.id JOIN entities lo ON p.legal_owner=lo.id `

const createProductStatement = `INSERT INTO products (id, legal_owner, display_name, summary, support_phone, support_email, homepage, logo_url, repo_url, issues_url, ontology) SELECT ?,lo.id,?,?,?,?,?,?,?,?,? FROM entities lo WHERE lo.pub_id=?`
func CreateProduct(p *Product, ctx context.Context) (*Product, rest.RestError) {
  txn, err := sqldb.DB.Begin()
  if err != nil {
    defer txn.Rollback()
    return nil, rest.ServerError("Could not create product record. (txn error)", err)
  }
  newP, restErr := CreateProductInTxn(p, ctx, txn)
  // txn already rolled back if in error, so we only need to commit if no error
  if err == nil {
    defer txn.Commit()
  }
  return newP, restErr
}

func CreateProductInTxn(p *Product, ctx context.Context, txn *sql.Tx) (*Product, rest.RestError) {
  var err error
  newId, restErr := entities.CreateEntityInTxn(txn)
  if restErr != nil {
    defer txn.Rollback()
		return nil, restErr
  }

  p.Id = nulls.NewInt64(newId)

	_, err = txn.Stmt(createProductQuery).Exec(newId, p.DisplayName, p.Summary, p.SupportPhone, p.SupportEmail, p.Homepage, p.LogoURL, p.RepoURL, p.IssuesURL, p.Ontology, p.LegalOwnerPubID)
	if err != nil {
    // TODO: can we do more to tell the cause of the failure? We assume it's due to malformed data with the HTTP code
    defer txn.Rollback()
    log.Print(err)
		return nil, rest.UnprocessableEntityError("Failure creating product.", err)
	}

  newProduct, err := GetProductByIDInTxn(p.Id.Int64, ctx, txn)
  if err != nil {
    return nil, rest.ServerError("Problem retrieving newly updated product.", err)
  }

  return newProduct, nil
}

const CommonProductGet string = `SELECT ` + CommonProductFields + CommonProductsFrom
const getProductStatement string = CommonProductGet + `WHERE e.pub_id=? `

// GetProduct retrieves a Product from a public ID string (UUID). Attempting to
// retrieve a non-existent Product results in a rest.NotFoundError. This is used
// primarily to retrieve a Product in response to an API request.
//
// Consider using GetProductByID to retrieve a Product from another backend/DB
// function. TODO: reference discussion of internal vs public IDs.
func GetProduct(pubId string, ctx context.Context) (*Product, rest.RestError) {
  return getProductHelper(getProductQuery, pubId, ctx, nil)
}

// GetProductInTxn retrieves a Product by public ID string (UUID) in the context
// of an existing transaction. See GetProduct.
func GetProductInTxn(pubId string, ctx context.Context, txn *sql.Tx) (*Product, rest.RestError) {
  return getProductHelper(getProductQuery, pubId, ctx, txn)
}

const getProductByIdStatement string = CommonProductGet + ` WHERE p.id=? `
// GetProductByID retrieves a Product by internal ID. As the internal ID must
// never be exposed to users, this method is exclusively for internal/backend
// use. Specifically, since Products are associated with other Entities through
// the internal ID (i.e., foreign keys use the internal ID), this function is
// most often used to retrieve a Product which is to be bundled in a response.
//
// Use GetProduct to retrieve a Product in response to an API request. TODO:
// reference discussion of internal vs public IDs.
func GetProductByID(id int64, ctx context.Context) (*Product, rest.RestError) {
  return getProductHelper(getProductByIdQuery, id, ctx, nil)
}

// GetProductByIDInTxn retrieves a Product by internal ID in the context of an
// existing transaction. See GetProductByID.
func GetProductByIDInTxn(id int64, ctx context.Context, txn *sql.Tx) (*Product, rest.RestError) {
  return getProductHelper(getProductByIdQuery, id, ctx, txn)
}

func getProductHelper(stmt *sql.Stmt, id interface{}, ctx context.Context, txn *sql.Tx) (*Product, rest.RestError) {
  if txn != nil {
    stmt = txn.Stmt(stmt)
  }
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, rest.ServerError("Error retrieving product.", err)
	}
	defer rows.Close()

	var product *Product
	for rows.Next() {
    var err error
    // The way the scanner works, it processes all the data each time. :(
    // 'product' gets updated with an equivalent structure while we gather up
    // the addresses.
    if product, err = ScanProduct(rows); err != nil {
      return nil, rest.ServerError(fmt.Sprintf("Problem getting data for product: '%v'", id), err)
    }
	} // TODO: we expect a single row
  if product != nil {
    product.FormatOut()
    return product, nil
  } else {
    return nil, rest.NotFoundError(fmt.Sprintf(`Product '%s' not found.`, id), nil)
  }
}

// BUG(zane@liquid-labs.com): UpdateProduct should use internal IDs if available
// on the Product struct. (I'm assuming this is slightly more efficient, though
// we should test.)

// UpdatesProduct updates the canonical Product record. Attempting to update a
// non-existent Product results in a rest.NotFoundError.
func UpdateProduct(p *Product, ctx context.Context) (*Product, rest.RestError) {
  txn, err := sqldb.DB.Begin()
  if err != nil {
    defer txn.Rollback()
    return nil, rest.ServerError("Could not update product record.", err)
  }

  newP, restErr := UpdateProductInTxn(p, ctx, txn)
  // txn already rolled back if in error, so we only need to commit if no error
  if restErr == nil {
    defer txn.Commit()
  }

  return newP, restErr
}

// UpdatesProductInTxn updates the canonical Product record within an existing
// transaction. See UpdateProduct.
func UpdateProductInTxn(p *Product, ctx context.Context, txn *sql.Tx) (*Product, rest.RestError) {
  var err error
  var updateStmt *sql.Stmt = txn.Stmt(updateProductQuery)
  _, err = updateStmt.Exec(p.LegalOwnerPubID, p.DisplayName, p.Summary, p.SupportPhone, p.SupportEmail, p.Homepage, p.LogoURL, p.RepoURL, p.IssuesURL, p.Ontology, p.PubId)
  if err != nil {
    if txn != nil {
      defer txn.Rollback()
    }
    return nil, rest.ServerError("Could not update product record.", err)
  }

  newProduct, err := GetProductInTxn(p.PubId.String, ctx, txn)
  if err != nil {
    return nil, rest.ServerError("Problem retrieving newly updated product.", err)
  }

  return newProduct, nil
}

// TODO: enable update of AuthID
const updateProductStatement = `UPDATE products p JOIN entities e ON p.id=e.id JOIN entities lo ON p.legal_owner=lo.id AND lo.pub_id=? SET p.legal_owner=lo.id, p.display_name=?, p.summary=?, p.support_phone=?, p.support_email=?, p.homepage=?, p.logo_url=?, p.repo_url=?, p.issues_url=?, p.ontology=?, e.last_updated=0 WHERE e.pub_id=?`
var createProductQuery, updateProductQuery, getProductQuery, getProductByAuthIdQuery, getProductByIdQuery *sql.Stmt
func SetupDB(db *sql.DB) {
  var err error
  if createProductQuery, err = db.Prepare(createProductStatement); err != nil {
    log.Fatalf("mysql: prepare create product stmt:\n%v\n%s", err, createProductStatement)
  }
  if getProductQuery, err = db.Prepare(getProductStatement); err != nil {
    log.Fatalf("mysql: prepare get product stmt:\n%v\nQuery: %s", err, getProductStatement)
  }
  if getProductByIdQuery, err = db.Prepare(getProductByIdStatement); err != nil {
    log.Fatalf("mysql: prepare get product by ID stmt:\n%v\n%s", err, getProductByIdStatement)
  }
  if updateProductQuery, err = db.Prepare(updateProductStatement); err != nil {
    log.Fatalf("mysql: prepare update product stmt:\n%v\n%s", err, updateProductStatement)
  }
}
