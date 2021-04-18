package product

import (
	"github.com/xeynse/XeynseJar_analytics/internal/config"
	productentity "github.com/xeynse/XeynseJar_analytics/internal/entity/product"
	db "github.com/xeynse/XeynseJar_analytics/internal/resource/db/dynamo"
)

type resource struct {
	config *config.Config
	db     db.Resource
}

type Resource interface {
	GetProductMacroDetails(productID string) (*productentity.ProductMacros, error)
}

func New(config *config.Config, dbRes db.Resource) Resource {
	return &resource{
		config: config,
		db:     dbRes,
	}
}

func (r resource) GetProductMacroDetails(homeID string) (*productentity.ProductMacros, error) {
	response := new(productentity.ProductMacros)
	err := r.db.GetItem(r.config.DataBase.Products.Table, "ProductID", homeID, response)
	if err != nil {
		return response, err
	}
	return response, nil
}
