package tools

import "time"

func TimeDurationReserveUser() time.Duration {
	// time.Now().Local().UTC() returns the time an hour before
	t_2 := time.Now().Local().UTC().Add(time.Hour * 12)
	t_1 := time.Now().Local().UTC()
	duration := t_2.Sub(t_1)
	return duration
}

func TimeReserveUser() time.Time {
	// time.Now().Local().UTC() returns the time an hour before
	return time.Now().Local().UTC().Add(time.Hour * 12)
}

func TimeDurationMessage() time.Duration {
	// time.Now().Local().UTC() returns the time an hour before
	t_2 := time.Now().Local().UTC().Add(time.Hour * 10)
	t_1 := time.Now().Local().UTC()
	duration := t_2.Sub(t_1)
	return duration
}

func TimeMessage() time.Time {
	// time.Now().Local().UTC() returns the time an hour before
	return time.Now().Local().UTC().Add(time.Hour * 10)
}

func GetTimeInt(date time.Time) int {
	// Calculate the duration
	duration := date.Sub(time.Now().Add(time.Hour))

	// Convert duration to int representing hours
	return int(duration.Hours())
}

func TimeToMicroseconds(t time.Time) int64 {
	return t.UnixNano() / 1000
}
