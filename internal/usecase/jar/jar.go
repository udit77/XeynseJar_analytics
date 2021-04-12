package jar

import (
	"errors"
	"sort"
	"time"

	"github.com/xeynse/XeynseJar_analytics/internal/config"
	"github.com/xeynse/XeynseJar_analytics/internal/entity"
	analyticsentity "github.com/xeynse/XeynseJar_analytics/internal/entity/analytics"
	"github.com/xeynse/XeynseJar_analytics/internal/repo/jar"
	"github.com/xeynse/XeynseJar_analytics/internal/util/common"
)

const (
	ContextDay   = "DAY"
	ContextWeek  = "WEEK"
	ContextMonth = "MONTH"
)

type resource struct {
	config *config.Config
	jar    jar.Repo
}

type Resource interface {
	GetAllJarStats(homeID string) ([]*analyticsentity.CurrentWeightStatus, error)
	GetJarStatByJarID(status *entity.JarStatusRequest) (*analyticsentity.CurrentWeightStatus, error)
	GetJarCalorieConsumption(status *entity.ConsumptionRequest) (*entity.JarConsumption, error)
}

func New(config *config.Config, jarRepo jar.Repo) Resource {
	return &resource{
		config: config,
		jar:    jarRepo,
	}
}

func (r *resource) GetAllJarStats(homeID string) ([]*analyticsentity.CurrentWeightStatus, error) {
	response, err := r.jar.GetAllJarStats(homeID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *resource) GetJarStatByJarID(status *entity.JarStatusRequest) (*analyticsentity.CurrentWeightStatus, error) {
	response, err := r.jar.GetJarStatByJarID(status.HomeID, status.JarID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *resource) GetJarCalorieConsumption(status *entity.ConsumptionRequest) (*entity.JarConsumption, error) {
	if status.Context == ContextDay {
		return r.GetCalorieConsumptionForDay(status)
	} else if status.Context == ContextWeek {
		return r.GetCalorieConsumptionForWeek(status)
	} else if status.Context == ContextMonth {
		return r.GetCalorieConsumptionForMonth(status)
	} else {
		return nil, errors.New("invalid request")
	}
}

func (r *resource) GetCalorieConsumptionForDay(status *entity.ConsumptionRequest) (*entity.JarConsumption, error) {
	var totalDayConsumption float64 = 0
	startTime := common.GetStartTimeForDay(time.Now())
	endTime := common.GetEndTimeForDay(time.Now())
	hourWiseDailyConsumption := make(map[int]float64)
	consumption, err := r.jar.GetCalorieConsumptionForDay(status.HomeID, status.JarID)
	if err != nil {
		return nil, err
	}
	for i := range consumption {
		weightDiff := consumption[i].WeightDiff
		updateTime := consumption[i].UpdateTime

		updateTimeUnix := updateTime.Unix()

		if updateTimeUnix >= startTime && updateTimeUnix <= endTime {
			hour := common.GetTimeByLocation(updateTime).Hour()
			if _, found := hourWiseDailyConsumption[hour]; !found {
				hourWiseDailyConsumption[hour] = 0
			}
			totalDayConsumption += weightDiff
			hourWiseDailyConsumption[hour] += weightDiff
		}
	}
	data := make([]*entity.ConsumptionForday, 0)
	for key, value := range hourWiseDailyConsumption {
		data = append(data, &entity.ConsumptionForday{
			Hour:  key,
			Value: value,
		})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Hour < data[j].Hour
	})
	return &entity.JarConsumption{
		TotalConsumption: totalDayConsumption,
		Data:             data,
	}, nil
}

func (r *resource) GetCalorieConsumptionForWeek(status *entity.ConsumptionRequest) (*entity.JarConsumption, error) {
	var totalWeekConsumption float64 = 0
	weekDayWiseConsumption := make(map[string]float64)
	consumption, err := r.jar.GetCalorieConsumptionForWeek(status.HomeID, status.JarID)
	if err != nil {
		return nil, err
	}
	for i := range consumption {
		weightDiff := consumption[i].WeightDiff
		updateTime := consumption[i].UpdateTime

		day := updateTime.Format("2006-01-02")

		if _, found := weekDayWiseConsumption[day]; !found {
			weekDayWiseConsumption[day] = 0
		}
		totalWeekConsumption += weightDiff
		weekDayWiseConsumption[day] += weightDiff
	}
	data := make([]*entity.ConsumptionForWeek, 0)
	for key, value := range weekDayWiseConsumption {
		day, _ := time.Parse("2006-01-02", key)
		data = append(data, &entity.ConsumptionForWeek{
			Date:  key,
			Day:   day.Weekday().String(),
			Value: value,
		})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Date < data[j].Date
	})
	return &entity.JarConsumption{
		TotalConsumption: totalWeekConsumption,
		Data:             data,
	}, nil
}

func (r *resource) GetCalorieConsumptionForMonth(status *entity.ConsumptionRequest) (*entity.JarConsumption, error) {
	var totalMonthConsumption float64 = 0
	dayWiseMonthlyConsumption := make(map[string]float64)
	consumption, err := r.jar.GetCalorieConsumptionForMonth(status.HomeID, status.JarID)
	if err != nil {
		return nil, err
	}
	for i := range consumption {
		weightDiff := consumption[i].WeightDiff
		updateTime := consumption[i].UpdateTime

		day := updateTime.Format("2006-01-02")

		if _, found := dayWiseMonthlyConsumption[day]; !found {
			dayWiseMonthlyConsumption[day] = 0
		}
		totalMonthConsumption += weightDiff
		dayWiseMonthlyConsumption[day] += weightDiff
	}
	data := make([]*entity.ConsumptionForMonth, 0)
	for key, value := range dayWiseMonthlyConsumption {
		data = append(data, &entity.ConsumptionForMonth{
			Date:  key,
			Value: value,
		})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Date < data[j].Date
	})
	return &entity.JarConsumption{
		TotalConsumption: totalMonthConsumption,
		Data:             data,
	}, nil
}
