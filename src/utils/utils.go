package uitls

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
	. "wanbu_stat_zzmm/src/logs"
)

var (
	Total int32
	Deal  int32
)

func CheckError(err error) {
	if err != nil {
		Logger.Critical(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func AddTotal() {
	atomic.AddInt32(&Total, 1)
}

func AddDeal() {
	atomic.AddInt32(&Deal, 1)
}

func Watching() {
	ti := time.Tick(time.Second * time.Duration(2))
	for {
		select {
		case <-ti:
			Logger.Info("Status : ", Deal, "/", Total)
		}
	}
}
