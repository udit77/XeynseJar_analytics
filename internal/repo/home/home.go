package home

import (
	"github.com/xeynse/XeynseJar_analytics/internal/config"
	"github.com/xeynse/XeynseJar_analytics/internal/entity"
	"github.com/xeynse/XeynseJar_analytics/internal/entity/home"
	db "github.com/xeynse/XeynseJar_analytics/internal/resource/db/dynamo"
)

type resource struct {
	config *config.Config
	db     db.Resource
}

type Resource interface {
	GetHomeConfiguration(homeID string) (home.Configuration, error)
}

func New(config *config.Config, dbRes db.Resource) Resource {
	return &resource{
		config: config,
		db:     dbRes,
	}
}

func (r resource) GetHomeConfiguration(homeID string) (home.Configuration, error) {
	homeConfig := new(entity.HomeConfigurationResponse)
	err := r.db.GetItem(r.config.DataBase.HomeConfiguration.Table, "HomeID", homeID, homeConfig)
	if err != nil {
		return home.Configuration{}, err
	}
	return homeConfig.Configuration, nil
}
