// @title LanLookout API
// @version 0.1
// @description Home network monitor: devices, internet health, security and DNS
// @contact.url   https://github.com/sarat03/lanlookout
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
// @BasePath /api/

package main

import (
	"flag"
	"os"
	// "net/http"

	// _ "net/http/pprof"

	// Import Swagger docs
	_ "github.com/sarat03/lanlookout/docs"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/dnslog"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/identify"
	"github.com/sarat03/lanlookout/internal/netcheck"
	"github.com/sarat03/lanlookout/internal/routines"
	"github.com/sarat03/lanlookout/internal/secscan"
	"github.com/sarat03/lanlookout/internal/web"
)

const dirPath = "/data/lanlookout"

// oldDirPath - data dir of WatchYourLAN, which LanLookout grew out of
const oldDirPath = "/data/WatchYourLAN"
const nodePath = ""

func main() {
	dirPtr := flag.String("d", dirPath, "Path to config dir")
	nodePtr := flag.String("n", nodePath, "Path to node modules")
	flag.Parse()

	// pprof - memory leak detect
	// go tool pprof -alloc_space http://localhost:8085/debug/pprof/heap
	// (pprof) web
	// (pprof) list db.Select
	//
	// go func() {
	// 	http.ListenAndServe("localhost:8085", nil)
	// }()

	// Existing WatchYourLAN installs keep their data: if only the old dir exists, use it.
	dir := *dirPtr
	if dir == dirPath && !exists(dirPath) && exists(oldDirPath) {
		dir = oldDirPath
	}

	// Make AppConfig
	conf.Start(dir, *nodePtr)

	gdb.Start()

	if conf.AppConfig.IdentifyEnable {
		identify.Start()
	}
	netcheck.Start()
	secscan.Start()
	dnslog.Start()
	routines.ScanRestart()
	routines.HistoryTrim()

	web.Gui()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
