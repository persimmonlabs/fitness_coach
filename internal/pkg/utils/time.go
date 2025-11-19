package utils

import (
	"fmt"
	"time"
)

// TimeFormat constants for common time formats
const (
	DateFormat         = "2006-01-02"
	TimeFormat         = "15:04:05"
	DateTimeFormat     = "2006-01-02 15:04:05"
	ISO8601Format      = "2006-01-02T15:04:05Z07:00"
	RFC3339Format      = time.RFC3339
	RFC3339NanoFormat  = time.RFC3339Nano
	CompactDateFormat  = "20060102"
	CompactTimeFormat  = "150405"
	HumanReadableDate  = "January 2, 2006"
	HumanReadableTime  = "3:04 PM"
	HumanReadableFull  = "January 2, 2006 at 3:04 PM"
)

// Common timezone locations
var (
	UTC        = time.UTC
	LocalTZ    *time.Location
	CommonTZs  map[string]*time.Location
)

func init() {
	LocalTZ = time.Local

	// Initialize common timezones
	CommonTZs = make(map[string]*time.Location)
	tzNames := []string{
		"America/New_York",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"Europe/London",
		"Europe/Paris",
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Australia/Sydney",
	}

	for _, tzName := range tzNames {
		if loc, err := time.LoadLocation(tzName); err == nil {
			CommonTZs[tzName] = loc
		}
	}
}

// NowUTC returns the current time in UTC
func NowUTC() time.Time {
	return time.Now().UTC()
}

// StartOfDay returns the start of the day (midnight) for the given time
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day (23:59:59.999999999) for the given time
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek returns the start of the week (Monday at midnight) for the given time
func StartOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	offset := int(weekday - time.Monday)
	if offset < 0 {
		offset += 7
	}
	return StartOfDay(t.AddDate(0, 0, -offset))
}

// EndOfWeek returns the end of the week (Sunday at 23:59:59.999999999) for the given time
func EndOfWeek(t time.Time) time.Time {
	return EndOfDay(StartOfWeek(t).AddDate(0, 0, 6))
}

// StartOfMonth returns the start of the month for the given time
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of the month for the given time
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// StartOfYear returns the start of the year for the given time
func StartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the end of the year for the given time
func EndOfYear(t time.Time) time.Time {
	return StartOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// ConvertToTimezone converts a time to a specific timezone
func ConvertToTimezone(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %w", err)
	}
	return t.In(loc), nil
}

// ConvertToUTC converts a time to UTC
func ConvertToUTC(t time.Time) time.Time {
	return t.UTC()
}

// ParseDate parses a date string in the format YYYY-MM-DD
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(DateFormat, dateStr)
}

// ParseDateTime parses a datetime string in the format YYYY-MM-DD HH:MM:SS
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse(DateTimeFormat, dateTimeStr)
}

// ParseISO8601 parses an ISO8601 formatted time string
func ParseISO8601(dateTimeStr string) (time.Time, error) {
	return time.Parse(ISO8601Format, dateTimeStr)
}

// ParseRFC3339 parses an RFC3339 formatted time string
func ParseRFC3339(dateTimeStr string) (time.Time, error) {
	return time.Parse(RFC3339Format, dateTimeStr)
}

// FormatDate formats a time as YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// FormatDateTime formats a time as YYYY-MM-DD HH:MM:SS
func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

// FormatISO8601 formats a time as ISO8601
func FormatISO8601(t time.Time) string {
	return t.Format(ISO8601Format)
}

// FormatRFC3339 formats a time as RFC3339
func FormatRFC3339(t time.Time) string {
	return t.Format(RFC3339Format)
}

// FormatHumanReadable formats a time in a human-readable format
func FormatHumanReadable(t time.Time) string {
	return t.Format(HumanReadableFull)
}

// DaysBetween calculates the number of days between two dates
func DaysBetween(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration.Hours() / 24)
}

// HoursBetween calculates the number of hours between two times
func HoursBetween(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration.Hours())
}

// IsToday checks if a time is today
func IsToday(t time.Time) bool {
	now := time.Now()
	return StartOfDay(t).Equal(StartOfDay(now))
}

// IsYesterday checks if a time is yesterday
func IsYesterday(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return StartOfDay(t).Equal(StartOfDay(yesterday))
}

// IsTomorrow checks if a time is tomorrow
func IsTomorrow(t time.Time) bool {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return StartOfDay(t).Equal(StartOfDay(tomorrow))
}

// IsWeekend checks if a time falls on a weekend (Saturday or Sunday)
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday checks if a time falls on a weekday (Monday-Friday)
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// AddBusinessDays adds a number of business days (skipping weekends)
func AddBusinessDays(t time.Time, days int) time.Time {
	current := t
	remaining := days

	if days > 0 {
		for remaining > 0 {
			current = current.AddDate(0, 0, 1)
			if IsWeekday(current) {
				remaining--
			}
		}
	} else if days < 0 {
		for remaining < 0 {
			current = current.AddDate(0, 0, -1)
			if IsWeekday(current) {
				remaining++
			}
		}
	}

	return current
}

// TimeAgo returns a human-readable string describing how long ago a time was
func TimeAgo(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if duration < 30*24*time.Hour {
		weeks := int(duration.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if duration < 365*24*time.Hour {
		months := int(duration.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(duration.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

// GetAge calculates age from a birthdate
func GetAge(birthdate time.Time) int {
	now := time.Now()
	age := now.Year() - birthdate.Year()

	// Adjust if birthday hasn't occurred this year yet
	if now.Month() < birthdate.Month() ||
		(now.Month() == birthdate.Month() && now.Day() < birthdate.Day()) {
		age--
	}

	return age
}

// IsExpired checks if a time is in the past
func IsExpired(t time.Time) bool {
	return t.Before(time.Now())
}

// TimeUntil returns a human-readable string describing time until a future time
func TimeUntil(t time.Time) string {
	now := time.Now()
	if t.Before(now) {
		return "expired"
	}

	duration := t.Sub(now)

	if duration < time.Minute {
		return "less than a minute"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}
}
