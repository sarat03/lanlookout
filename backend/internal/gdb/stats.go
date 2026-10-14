package gdb

import (
	"github.com/sarat03/lanlookout/internal/check"
	"github.com/sarat03/lanlookout/internal/models"
)

// AddNetStat - store one internet check
func AddNetStat(s models.NetStat) {
	check.IfError(db.Table("netstats").Create(&s).Error)
}

// SelectNetStats - checks since date, oldest first
func SelectNetStats(since string) (stats []models.NetStat) {
	db.Table("netstats").Where("\"DATE\" >= ?", since).Order("\"ID\"").Find(&stats)
	return stats
}

// AddDNSCounts - add query counts for one day; rows are keyed by day+ip+domain
func AddDNSCounts(day string, counts map[[2]string]models.DNSStat) {

	for key, c := range counts {
		var row models.DNSStat
		db.Table("dnsstats").Where("\"DAY\" = ? AND \"IP\" = ? AND \"DOMAIN\" = ?", day, key[0], key[1]).Limit(1).Find(&row)

		row.Day, row.IP, row.Domain = day, key[0], key[1]
		row.Count += c.Count
		row.Blocked += c.Blocked
		row.Threat = c.Threat
		check.IfError(db.Table("dnsstats").Save(&row).Error)
	}
}

// SelectDNSTop - most queried domains of one IP since day
func SelectDNSTop(ip, since string, limit int) (stats []models.DNSStat) {
	db.Table("dnsstats").
		Select("\"DOMAIN\", SUM(\"COUNT\") AS \"COUNT\", SUM(\"BLOCKED\") AS \"BLOCKED\", MAX(\"THREAT\") AS \"THREAT\"").
		Where("\"IP\" = ? AND \"DAY\" >= ?", ip, since).
		Group("DOMAIN").Order("\"COUNT\" DESC").Limit(limit).
		Find(&stats)
	return stats
}

// DeleteOldStats - trim netstats and dnsstats
func DeleteOldStats(date string) {
	check.IfError(db.Table("netstats").Where("\"DATE\" < ?", date).Delete(&models.NetStat{}).Error)
	check.IfError(db.Table("dnsstats").Where("\"DAY\" < ?", date[:10]).Delete(&models.DNSStat{}).Error)
}
