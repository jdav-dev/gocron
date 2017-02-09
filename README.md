Golang CRON Expression Generator
================================
[![Build Status](https://travis-ci.org/jdav-dev/gocron.svg?branch=master)](https://travis-ci.org/jdav-dev/gocron)

Takes a `time.Duration` and returns the closest CRON expression to fit that
interval.  Can also offset that CRON expression with by a second
`time.Duration`.  Returned expressions fit the standard/default CRON
implementation.

Usage
-----
Import:
```golang
import (
	"time"

	"github.com/jdav-dev/cron"
)
```

Every minute (minimum interval):
```golang
s, _ := gocron.IntervalToSchedule(time.Minute)
s.Expression()
```
returns `* * * * *`

Every 15 minutes:
```golang
s, _ := gocron.IntervalToSchedule(15 * time.Minute)
s.Expression()
```
returns `0,15,30,45 * * * *`

Once a year (maximum interval):
```golang
s, _ := gocron.IntervalToSchedule(365 * 24 * time.Hour)
s.Expression()
```
returns `0 0 1 1 *`

Every 15 minutes, offset by 5 minutes:
```golang
s, _ := gocron.OffsetIntervalToSchedule(15 * time.Minute, 5 * time.Minute)
s.Expression()
```
returns `5,20,35,50 * * * *`
