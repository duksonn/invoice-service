package utils

import (
	"time"
)

var Now func() time.Time = time.Now

func PatchNow(now func() time.Time) {
	Now = now
}

func RestoreNow() {
	Now = time.Now
}
