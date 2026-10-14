package api

import (
	"strconv"

	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/models"
)

func getHostByID(idStr string) (oneHost models.Host) {

	id, _ := strconv.Atoi(idStr)
	oneHost = gdb.SelectByID(id)

	return oneHost
}
