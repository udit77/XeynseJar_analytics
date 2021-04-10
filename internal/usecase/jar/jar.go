package jar

import (
	"github.com/xeynse/XeynseJar_analytics/internal/config"
	analyticsentity "github.com/xeynse/XeynseJar_analytics/internal/entity/analytics"
	"github.com/xeynse/XeynseJar_analytics/internal/repo/jar"
)

type resource struct {
	config *config.Config
	jar    jar.Repo
}

type Resource interface {
	GetAllJarStats(homeID string) ([]*analyticsentity.CurrentWeightStatus, error)
	GetJarStatByJarID(homeID string, jarID string) (*analyticsentity.CurrentWeightStatus, error)
}

func New(config *config.Config, jarRepo jar.Repo) Resource {
	return &resource{
		config: config,
	}
}

func (r *resource) GetAllJarStats(homeID string) ([]*analyticsentity.CurrentWeightStatus, error) {
	response, err := r.jar.GetAllJarStats(homeID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *resource) GetJarStatByJarID(homeID string, jarID string) (*analyticsentity.CurrentWeightStatus, error) {
	response, err := r.jar.GetJarStatByJarID(homeID, jarID)
	if err != nil {
		return nil, err
	}
	return response, nil
}
