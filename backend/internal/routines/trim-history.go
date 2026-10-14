package routines

import (
	"log/slog"
	"time"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
)

// HistoryTrim - routine for History
func HistoryTrim() {

	go func() {
		for {
			time.Sleep(time.Duration(1) * time.Hour) // Every hour

			hours := conf.AppConfig.TrimHist
			nowMinus := time.Now().Add(-time.Duration(hours) * time.Hour)
			date := nowMinus.Format("2006-01-02 15:04:05")

			slog.Info("Removing all History before", "date", date)

			n := gdb.DeleteOldHistory(date)
			slog.Info("Removed records from History", "n", n)

			// ponytail: events kept 7x longer than history, they are tiny
			evDate := time.Now().Add(-7 * time.Duration(hours) * time.Hour).Format("2006-01-02 15:04:05")
			gdb.DeleteOldEvents(evDate)
			gdb.DeleteOldStats(evDate)
		}
	}()
}
