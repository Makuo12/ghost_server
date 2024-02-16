package tools

import "time"

type Weekend struct {
	Friday   time.Time
	Saturday time.Time
	Sunday   time.Time
}

func ListWeekends() []Weekend {
	var weekends []Weekend

	// Get the current time
	now := time.Now()

	// Find the next Friday from the current day
	for now.Weekday() != time.Friday {
		now = now.AddDate(0, 0, 1)
	}

	// Iterate through each week until the next year
	for currentWeek := 0; currentWeek <= 51; currentWeek++ {
		// Find the first day of the week
		weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)

		// Calculate the days to add to get to the weekend
		daysToAdd := 5 - int(weekStart.Weekday())

		// Calculate the weekend dates
		friday := weekStart.AddDate(0, 0, daysToAdd)
		saturday := friday.AddDate(0, 0, 1)
		sunday := friday.AddDate(0, 0, 2)

		// Create a Weekend struct and append it to the list
		weekend := Weekend{
			Friday:   friday,
			Saturday: saturday,
			Sunday:   sunday,
		}
		weekends = append(weekends, weekend)

		// Move to the next week
		now = now.AddDate(0, 0, 7)
	}

	return weekends
}

type Week struct {
	Monday time.Time
	Sunday time.Time
}

// ListWeeks generates a list of weeks from the current week to the next year
func ListWeeks() []Week {
	var weeks []Week

	// Get the current time
	now := time.Now()

	// Iterate through each week until the next year
	for currentWeek := 0; currentWeek <= 51; currentWeek++ {
		// Find the first day of the week (Monday)
		weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)

		// Calculate the days to add to get to the end of the week (Sunday)
		daysToAdd := 6 - int(weekStart.Weekday())

		// Calculate the Monday and Sunday dates
		monday := weekStart
		sunday := weekStart.AddDate(0, 0, daysToAdd)

		// Create a Week struct and append it to the list
		week := Week{
			Monday: monday,
			Sunday: sunday,
		}
		weeks = append(weeks, week)

		// Move to the next week
		now = now.AddDate(0, 0, 7)
	}

	return weeks
}

// Month struct to represent the start and end dates of a month
type Month struct {
	StartDateOfMonth time.Time
	EndDateOfMonth   time.Time
}

// ListMonths generates a list of months for the next two years
func ListMonths() []Month {
	var months []Month

	// Get the current time
	now := time.Now()

	// Iterate through each month for the next two years
	for year := now.Year(); year <= now.Year()+1; year++ {
		for currentMonth := time.January; currentMonth <= time.December; currentMonth++ {
			// Find the first day of the month
			monthStart := time.Date(year, currentMonth, 1, 0, 0, 0, 0, now.Location())

			// Find the last day of the month
			lastDayOfMonth := time.Date(year, currentMonth, 1, 0, 0, 0, 0, now.Location()).
				AddDate(0, 1, -1)

			// Create a Month struct and append it to the list
			month := Month{
				StartDateOfMonth: monthStart,
				EndDateOfMonth:   lastDayOfMonth,
			}
			months = append(months, month)
		}
	}

	return months
}

type ExDateAdd struct {
	StartDate time.Time
	EndDate   time.Time
}

func ListDateIntervals(space int) []ExDateAdd {
	var intervals []ExDateAdd

	// Get the current time
	now := time.Now()

	// Initialize the start date as the current date
	startDate := now

	// Iterate to generate date intervals until next year
	for startDate.Before(now.AddDate(1, 0, 0)) {
		// Calculate the end date with a 5-day difference
		endDate := startDate.AddDate(0, 0, space)

		// Create an ExDateAdd struct and append it to the list
		interval := ExDateAdd{
			StartDate: startDate,
			EndDate:   endDate,
		}
		intervals = append(intervals, interval)

		// Move to the next start date
		startDate = endDate.AddDate(0, 0, 1)
	}

	return intervals
}
