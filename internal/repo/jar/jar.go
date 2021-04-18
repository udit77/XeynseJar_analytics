package jar

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/xeynse/XeynseJar_analytics/internal/entity"
	analyticsentity "github.com/xeynse/XeynseJar_analytics/internal/entity/analytics"
	db "github.com/xeynse/XeynseJar_analytics/internal/resource/db/analytics"
	"github.com/xeynse/XeynseJar_analytics/internal/util/common"
)

type resource struct {
	dbResource db.Resource
	statements sqlStatement
}

type sqlStatement struct {
	fetchWeightStatByHomeID *sqlx.Stmt
	fetchWeightStatByJarID  *sqlx.Stmt
	fetchConsumption        *sqlx.Stmt
}

type Repo interface {
	InsertJarStateData(map[string]*entity.JarReported) error
	GetAllJarStats(homeID string) ([]*analyticsentity.CurrentWeightStatus, error)
	GetJarStatByJarID(homeID string, jarID string) (*analyticsentity.CurrentWeightStatus, error)
	GetCalorieConsumptionForDay(homeID string, jarID string) ([]*analyticsentity.Consumption, error)
	GetCalorieConsumptionForWeek(homeID string, jarID string) ([]*analyticsentity.Consumption, error)
	GetCalorieConsumptionForMonth(homeID string, jarID string) ([]*analyticsentity.Consumption, error)
}

func New(dbRes db.Resource) Repo {
	return &resource{
		dbResource: dbRes,
		statements: sqlStatement{
			fetchWeightStatByHomeID: dbRes.PreparexFatal(queryCurrentWeightByHomeID),
			fetchWeightStatByJarID:  dbRes.PreparexFatal(queryCurrentWeightByJarID),
			fetchConsumption:        dbRes.PreparexFatal(queryFetchConsumption),
		},
	}
}

func (r *resource) InsertJarStateData(jarState map[string]*entity.JarReported) error {
	if jarState == nil {
		return errors.New("[InsertJarStatus] empty jar state")
	}
	dbValues := make([]string, 0)
	for key, jarValue := range jarState {
		value := jarValue.Nodes.Node0
		valueString := fmt.Sprintf(` ('%v','%v',%v,%v,%v,%v,%v,'%v', now()) `, value.HomeID, key, value.WeightStart, value.WeightDiff, value.WeightCurrent, value.ZAxisG, value.Error, value.Type)
		dbValues = append(dbValues, valueString)
	}

	dbValueString := strings.Join(dbValues, ", ")
	query := fmt.Sprintf("%v%v", queryInsertJarStatusData, dbValueString)
	log.Println(query)
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

func (r *resource) GetCalorieConsumptionForDay(homeID string, jarID string) ([]*analyticsentity.Consumption, error) {
	response := make([]*analyticsentity.Consumption, 0)
	rows, err := r.statements.fetchConsumption.Queryx(homeID, jarID, common.GetTodayTimeStringByLocation())
	if err != nil {
		return response, err
	}
	defer rows.Close()
	for rows.Next() {
		analytics := &analyticsentity.Consumption{}
		err := rows.StructScan(analytics)
		if err != nil {
			log.Println(err)
			continue
		}
		response = append(response, analytics)
	}
	return response, nil
}

func (r *resource) GetCalorieConsumptionForWeek(homeID string, jarID string) ([]*analyticsentity.Consumption, error) {
	response := make([]*analyticsentity.Consumption, 0)
	rows, err := r.statements.fetchConsumption.Queryx(homeID, jarID, common.GetWeekTimeStringByLocation())
	if err != nil {
		return response, err
	}
	defer rows.Close()
	for rows.Next() {
		analytics := &analyticsentity.Consumption{}
		err := rows.StructScan(analytics)
		if err != nil {
			log.Println(err)
			continue
		}
		response = append(response, analytics)
	}
	return response, nil
}

func (r *resource) GetCalorieConsumptionForMonth(homeID string, jarID string) ([]*analyticsentity.Consumption, error) {
	response := make([]*analyticsentity.Consumption, 0)
	rows, err := r.statements.fetchConsumption.Queryx(homeID, jarID, common.GetMonthTimeStringByLocation())
	if err != nil {
		return response, err
	}
	defer rows.Close()
	for rows.Next() {
		analytics := &analyticsentity.Consumption{}
		err := rows.StructScan(analytics)
		if err != nil {
			log.Println(err)
			continue
		}
		response = append(response, analytics)
	}
	return response, nil
}
