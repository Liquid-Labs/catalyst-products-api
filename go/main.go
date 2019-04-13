package main

import (
  "github.com/Liquid-Labs/catalyst-core-api/go/restserv"
  // core resources
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/entities"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/users"

  "github.com/Liquid-Labs/catalyst-products-api/go/resources/products"
  "github.com/Liquid-Labs/go-api/sqldb"
)

func main() {
  sqldb.RegisterSetup(entities.SetupDB)
  sqldb.RegisterSetup(users.SetupDB)
  sqldb.RegisterSetup(products.SetupDB)
  sqldb.InitDB()
  restserv.RegisterResource(products.InitAPI)
  restserv.Init()
}
