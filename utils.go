package discordgo

import (
	"strconv"
	"time"
)

func ContainsIDObject(haystack []IDGettable, needle IDGettable) (contains bool) {
	if len(haystack) < 1 {
		return false
	}

	for _, item := range haystack {
		if item.GetID() == needle.GetID() {
			return true
		}
	}

	return false
}

func SnowflakeToTime(snowflake string) (returnTime time.Time, err error) {
	n, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return
	}

	timestamp := (n >> 22) + 1420070400000
	returnTime = time.Unix(timestamp, 0).UTC()
	return
}
