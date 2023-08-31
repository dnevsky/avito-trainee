package helpers

import (
	"regexp"
	"strconv"
	"time"
)

func ParseStringToTime(str string) (time.Time, error) {
	ttlTime := time.Now()

	re, err := regexp.Compile(`((\d+)(\w+))`)
	if err != nil {
		return ttlTime, err
	}

	res := re.FindAllStringSubmatch(str, -1)

	for _, v := range res {
		value, err := strconv.Atoi(v[2])
		if err != nil {
			return ttlTime, err
		}

		typ := v[3]

		switch typ {
		case "s":
			if value > 0 {
				ttlTime = ttlTime.Add(time.Second * time.Duration(value))
			}
		case "m":
			if value > 0 {
				ttlTime = ttlTime.Add(time.Minute * time.Duration(value))
			}
		case "h":
			if value > 0 {
				ttlTime = ttlTime.Add(time.Hour * time.Duration(value))
			}
		case "d":
			if value > 0 {
				ttlTime = ttlTime.AddDate(0, 0, value)
			}
		case "w":
			if value > 0 {
				ttlTime = ttlTime.Add(time.Hour * 24 * 7 * time.Duration(value))
			}
		case "y":
			if value > 0 {
				ttlTime = ttlTime.AddDate(value, 0, 0)
			}
		}
	}

	return ttlTime, nil
}
