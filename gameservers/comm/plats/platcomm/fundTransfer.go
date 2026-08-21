package platcomm

import (
	"log/slog"
	"net/url"
	"time"
)

type FnFundTransfer = func() bool

type FnCheckOrder = func() bool

func IsTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	urlerr, ok := err.(*url.Error)
	return ok && urlerr.Timeout()
}

func GetTransStatus(err error) string {
	if err != nil {
		if IsTimeoutErr(err) {
			return "TIMEOUT"
		}

		return "ERROR:" + err.Error()
	}
	return "SUCCESS"
}

func FundTransferWithRetry(transfer FnFundTransfer, check FnCheckOrder) (retry bool) {
	retry = transfer()

	for i := 0; i < 3 && retry; i++ {
		time.Sleep(time.Second * 5)
		retry = check()
	}

	return retry
}

func FundTransferWithRetryAndInterval(transfer FnFundTransfer, check FnCheckOrder, interval time.Duration) (retry bool) {
	retry = transfer()

	for i := 0; i < 3 && retry; i++ {
		time.Sleep(interval)
		retry = check()
	}

	return retry
}

var (
	sleepSecs = []time.Duration{2, 4, 8, 16, 32, 64, 128}
)

func ExcuteAndCheck(transfer FnFundTransfer, check FnCheckOrder) {
	needcheck := transfer()
	if !needcheck {
		return
	}

	for _, sec := range sleepSecs {
		retry := check()
		if !retry {
			break
		}
		time.Sleep(time.Second * sec)
	}
}

func Retry(maxcount int, sleep time.Duration, fn FnCheckOrder) {
	for i := 0; i < maxcount; i++ {
		retry := fn()
		if !retry {
			break
		}

		if i+1 == maxcount {
			break
		}
		time.Sleep(sleep)
		slog.Warn("Retry", "i", i+1, "maxcount", maxcount)
	}
}
