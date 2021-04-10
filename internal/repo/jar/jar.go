package jar

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xeynse/XeynseJar_analytics/internal/entity"
	analyticsentity "github.com/xeynse/XeynseJar_analytics/internal/entity/analytics"
	"github.com/xeynse/XeynseJar_analytics/internal/resource/db"
)

type resource struct {
	dbResource db.Resource
	statements sqlStatement
}

type sqlStatement struct {
	fetchWeightStatByHomeID *sqlx.Stmt
	fetchWeightStatByJarID  *sqlx.Stmt
}

type Repo interface {
	InsertJarStateData(map[string]entity.JarState) error
	GetAllJarStats(homeID string) ([]*analyticsentity.CurrentWeightStatus, error)
	GetJarStatByJarID(homeID string, jarID string) (*analyticsentity.CurrentWeightStatus, error)
}

func New(dbRes db.Resource) Repo {
	return &resource{
		dbResource: dbRes,
		statements: sqlStatement{
			fetchWeightStatByHomeID: dbRes.PreparexFatal(queryCurrentWeightByHomeID),
			fetchWeightStatByJarID:  dbRes.PreparexFatal(queryCurrentWeightByJarID),
		},
	}
}

func (r *resource) InsertJarStateData(state map[string]entity.JarState) error {
	dbValues := make([]string, 0)
	for key, value := range state {
		valueString := fmt.Sprintf(" (%v,%v,%v,%v,%v,%v,%v,%v,%v) ", value.HomeID, key, value.WeightStart, value.WeightDiff, value.WeightCurrent, value.ZAxisG, value.Error, value.Type, time.Now().UTC())
		dbValues = append(dbValues, valueString)
	}

	dbValueString := strings.Join(dbValues, ", ")
	query := fmt.Sprintf("%v%v", queryInsertJarStatusData, dbValueString)
	_, err := r.dbResource.GetClient().Exec(query)
	return err
}

func (r *resource) GetAllJarStats(homeID string) ([]*analyticsentity.CurrentWeightStatus, error) {
	response := make([]*analyticsentity.CurrentWeightStatus, 0)
	rows, err := r.statements.fetchWeightStatByHomeID.Queryx(homeID)
	if err != nil {
		return response, err
	}
	defer rows.Close()
	for rows.Next() {
		analytics := &analyticsentity.CurrentWeightStatus{}
		err := rows.StructScan(analytics)
		if err != nil {
			log.Println(err)
			continue
		}
		response = append(response, analytics)
	}
	return response, nil
}

func (r *resource) GetJarStatByJarID(homeID string, jarID string) (*analyticsentity.CurrentWeightStatus, error) {
	response := &analyticsentity.CurrentWeightStatus{}
	err := r.statements.fetchWeightStatByJarID.QueryRowx(homeID, jarID).StructScan(response)
	if err != nil {
		return response, err
	}
	return response, nil
}
