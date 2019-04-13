package products

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/Liquid-Labs/catalyst-core-api/go/handlers"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "/products is alive\n")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  var product *Product = &Product{}
  if _, restErr := handlers.CheckAndExtract(w, r, product, `Product`); restErr != nil {
    return // response handled by CheckAndExtract
  } else {
    handlers.DoCreate(w, r, CreateProduct, product, `Product`)
  }
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
  if _, restErr := handlers.BasicAuthCheck(w, r); restErr != nil {
    return // response handled by BasicAuthCheck
  } else {
    vars := mux.Vars(r)
    pubID := vars["pubId"]

    handlers.DoGetDetail(w, r, GetProduct, pubID, `Product`)
  }
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
  var newData *Product = &Product{}
  if _, restErr := handlers.CheckAndExtract(w, r, newData, `Product`); restErr != nil {
    return // response handled by CheckAndExtract
  } else {
    vars := mux.Vars(r)
    pubID := vars["pubId"]

    handlers.DoUpdate(w, r, UpdateProduct, newData, pubID, `Product`)
  }
}

const uuidRE = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}`

func InitAPI(r *mux.Router) {
  r.HandleFunc("/products/", pingHandler).Methods("PING")
  r.HandleFunc("/products/", createHandler).Methods("POST")
  r.HandleFunc("/products/{pubId:" + uuidRE + "}/", detailHandler).Methods("GET")
  r.HandleFunc("/products/{pubId:" + uuidRE + "}/", updateHandler).Methods("PUT")
}
