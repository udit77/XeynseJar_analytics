package jar

const (
	queryInsertJarStatusData = `
		INSERT INTO public.transactions_jar_status (home_id, jar_id, weight_start, weight_diff, weight_current, zaxis_g, error, type, update_time) VALUES
	`
	queryCurrentWeightByHomeID = `
		SELECT jar_id, weight_current from public.transactions_jar_status where home_id = $1 order by update_time desc limit 1
	`
	queryCurrentWeightByJarID = `
		SELECT jar_id, weight_current from public.transactions_jar_status where home_id = $1 and jar_id = $2 order by update_time desc limit 1
	`
	queryFetchConsumption = `
		SELECT jar_id, weight_diff, update_time from public.transactions_jar_status where home_id = $1 and jar_id = $2 and update_time >= $3
	`
)
