package gdb

import (
	"github.com/sarat03/lanlookout/internal/check"
	"github.com/sarat03/lanlookout/internal/models"
)

// Update - update or create host
func Update(table string, oneHost models.Host) {

	tab := db.Table(table)
	result := tab.Save(&oneHost)
	check.IfError(result.Error)
}

// Delete - delete host from DB
func Delete(table string, id int) {

	tab := db.Table(table)
	result := tab.Delete(&models.Host{}, id)
	check.IfError(result.Error)
}

// DeleteOldHistory - delete a list of hosts from History
func DeleteOldHistory(date string) int64 {

	tab := db.Table("history")
	result := tab.Where("\"DATE\" < ?", date).Delete(&models.Host{})
	check.IfError(result.Error)

	return result.RowsAffected
}

// Clear - delete all hosts from table
func Clear(table string) {

	tab := db.Table(table)
	result := tab.Where("1 = 1").Delete(&models.Host{})
	check.IfError(result.Error)
}

// DeleteHistoryByMAC - remove all history of one host
func DeleteHistoryByMAC(mac string) {

	tab := db.Table("history")
	result := tab.Where("\"MAC\" = ?", mac).Delete(&models.Host{})
	check.IfError(result.Error)
}

// AddEvent - store one alert
func AddEvent(e models.Event) {

	tab := db.Table("events")
	result := tab.Create(&e)
	check.IfError(result.Error)
}

// DeleteOldEvents - delete events older than date
func DeleteOldEvents(date string) int64 {

	tab := db.Table("events")
	result := tab.Where("\"DATE\" < ?", date).Delete(&models.Event{})
	check.IfError(result.Error)

	return result.RowsAffected
}

// UpdateScan - write only what a scan learns, so edits made meanwhile
// (name, watch, tags, ports from the security scan) aren't overwritten.
func UpdateScan(h models.Host, withName, withType bool) {

	cols := []string{"IFACE", "IP", "DATE", "NOW", "DNS", "MODEL", "OS", "SERVICES"}
	if withName {
		cols = append(cols, "NAME")
	}
	if withType {
		cols = append(cols, "TYPE")
	}
	result := db.Table("now").Where("\"ID\" = ?", h.ID).Select(cols).Updates(&h)
	check.IfError(result.Error)
}

// SetPorts - store the security scan result of one host
func SetPorts(id int, ports string) {
	result := db.Table("now").Where("\"ID\" = ?", id).Update("PORTS", ports)
	check.IfError(result.Error)
}
