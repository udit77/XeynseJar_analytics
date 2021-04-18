package jar

import (
	"errors"
	"sort"
	"time"

	"github.com/xeynse/XeynseJar_analytics/internal/config"
	"github.com/xeynse/XeynseJar_analytics/internal/entity"
	analyticsentity "github.com/xeynse/XeynseJar_analytics/internal/entity/analytics"
	"github.com/xeynse/XeynseJar_analytics/internal/repo/home"
	"github.com/xeynse/XeynseJar_analytics/internal/repo/jar"
	"github.com/xeynse/XeynseJar_analytics/internal/repo/product"

	"github.com/xeynse/XeynseJar_analytics/internal/util/common"
)

const (
	ContextDay   = "DAY"
	ContextWeek  = "WEEK"
	ContextMonth = "MONTH"
)

var WeekDay = [7]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

type resource struct {
	config     *config.Config
	jar        jar.Repo
	homeconfig home.Resource
	product    product.Resource
}

type Resource interface {
	GetAllJarStats(homeID string) ([]*analyticsentity.CurrentWeightStatus, error)
	GetJarStatByJarID(status *entity.JarStatusRequest) (*analyticsentity.CurrentWeightStatus, error)
	GetJarCalorieConsumption(status *entity.ConsumptionRequest) (*entity.JarConsumption, error)
	GetMacroNutrientConsumption(request *entity.ConsumptionRequest) (*entity.JarConsumption, error)
}

func New(config *config.Config, homeconfigRepo home.Resource, productRepo product.Resource, jarRepo jar.Repo) Resource {
	return &resource{
		config:     config,
		jar:        jarRepo,
		homeconfig: homeconfigRepo,
		product:    productRepo,
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
	consumption, err := r.jar.GetJarConsumptionForDay(status.HomeID, status.JarID)
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

	t1 := time.Now()
	weekDay := int(t1.Weekday())
	val := t1.Day() - (weekDay - 1)
	if weekDay == 0 {
		val = t1.Day() - 6
	}
	for i := t1.Day(); i >= val; i-- {
		date := time.Now().AddDate(0, 0, -1*(i-val)).Format("2006-01-02")
		weekDayWiseConsumption[date] = 0
	}

	consumption, err := r.jar.GetJarConsumptionForWeek(status.HomeID, status.JarID)
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
	dateToday := time.Now().Day()
	var totalMonthConsumption float64 = 0
	dayWiseMonthlyConsumption := make(map[string]float64)

	for i := dateToday - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -1*i).Format("2006-01-02")
		dayWiseMonthlyConsumption[date] = 0
	}

	consumption, err := r.jar.GetJarConsumptionForMonth(status.HomeID, status.JarID)
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

func (r *resource) GetMacroNutrientConsumption(status *entity.ConsumptionRequest) (*entity.JarConsumption, error) {
	if status.Context == ContextDay {
		return r.GetMacroConsumptionForDay(status)
	} else if status.Context == ContextWeek {
		return r.GetMacroConsumptionForWeek(status)
	} else if status.Context == ContextMonth {
		return r.GetMacroConsumptionForMonth(status)
	} else {
		return nil, errors.New("invalid request")
	}
}

func (r *resource) GetMacroConsumptionForDay(status *entity.ConsumptionRequest) (*entity.JarConsumption, error) {
	totalMacros := &entity.Macros{}
	var totalDayConsumption float64 = 0
	startTime := common.GetStartTimeForDay(time.Now())
	endTime := common.GetEndTimeForDay(time.Now())
	hourWiseDailyConsumption := make(map[int]float64)
	consumption, err := r.jar.GetJarConsumptionForDay(status.HomeID, status.JarID)
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

	productMacros, err := r.product.GetProductMacroDetails("101")
	if err == nil {
		totalMacros = &entity.Macros{
			Carbohydrates: (float64(productMacros.Carbohydrates) * totalDayConsumption) / 100,
			Fat:           (float64(productMacros.Fat) * totalDayConsumption) / 100,
			Protein:       (float64(productMacros.Protein) * totalDayConsumption) / 100,
		}
		for i := range data {
			data[i].Macros = &entity.Macros{
				Carbohydrates: (float64(productMacros.Carbohydrates) * data[i].Value) / 100,
				Fat:           (float64(productMacros.Fat) * data[i].Value) / 100,
				Protein:       (float64(productMacros.Protein) * data[i].Value) / 100,
			}
		}
	}

	return &entity.JarConsumption{
		TotalConsumption: totalDayConsumption,
		Macro:            totalMacros,
		Data:             data,
	}, nil
}

func (r *resource) GetMacroConsumptionForWeek(status *entity.ConsumptionRequest) (*entity.JarConsumption, error) {
	totalMacros := &entity.Macros{}
	var totalWeekConsumption float64 = 0
	weekDayWiseConsumption := make(map[string]float64)

	t1 := time.Now()
	weekDay := int(t1.Weekday())
	val := t1.Day() - (weekDay - 1)
	if weekDay == 0 {
		val = t1.Day() - 6
	}
	for i := t1.Day(); i >= val; i-- {
		date := time.Now().AddDate(0, 0, -1*(i-val)).Format("2006-01-02")
		weekDayWiseConsumption[date] = 0
	}

	consumption, err := r.jar.GetJarConsumptionForWeek(status.HomeID, status.JarID)
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

	productMacros, err := r.product.GetProductMacroDetails("101")
	if err == nil {
		totalMacros = &entity.Macros{
			Carbohydrates: (float64(productMacros.Carbohydrates) * totalWeekConsumption) / 100,
			Fat:           (float64(productMacros.Fat) * totalWeekConsumption) / 100,
			Protein:       (float64(productMacros.Protein) * totalWeekConsumption) / 100,
		}
		for i := range data {
			data[i].Macros = &entity.Macros{
				Carbohydrates: (float64(productMacros.Carbohydrates) * data[i].Value) / 100,
				Fat:           (float64(productMacros.Fat) * data[i].Value) / 100,
				Protein:       (float64(productMacros.Protein) * data[i].Value) / 100,
			}
		}
	}

	return &entity.JarConsumption{
		TotalConsumption: totalWeekConsumption,
		Macro:            totalMacros,
		Data:             data,
	}, nil
}

func (r *resource) GetMacroConsumptionForMonth(status *entity.ConsumptionRequest) (*entity.JarConsumption, error) {
	totalMacros := &entity.Macros{}
	dateToday := time.Now().Day()
	var totalMonthConsumption float64 = 0
	dayWiseMonthlyConsumption := make(map[string]float64)

	for i := dateToday - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -1*i).Format("2006-01-02")
		dayWiseMonthlyConsumption[date] = 0
	}

	consumption, err := r.jar.GetJarConsumptionForMonth(status.HomeID, status.JarID)
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

	productMacros, err := r.product.GetProductMacroDetails("101")
	if err == nil {
		totalMacros = &entity.Macros{
			Carbohydrates: (float64(productMacros.Carbohydrates) * totalMonthConsumption) / 100,
			Fat:           (float64(productMacros.Fat) * totalMonthConsumption) / 100,
			Protein:       (float64(productMacros.Protein) * totalMonthConsumption) / 100,
		}
		for i := range data {
			data[i].Macros = &entity.Macros{
				Carbohydrates: (float64(productMacros.Carbohydrates) * data[i].Value) / 100,
				Fat:           (float64(productMacros.Fat) * data[i].Value) / 100,
				Protein:       (float64(productMacros.Protein) * data[i].Value) / 100,
			}
		}
	}

	return &entity.JarConsumption{
		TotalConsumption: totalMonthConsumption,
		Macro:            totalMacros,
		Data:             data,
	}, nil
}
