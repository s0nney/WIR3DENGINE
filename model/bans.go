package model

import "WIR3DENGINE/utils"

func IsBanned(ip string) bool {
	var count int
	err := utils.Db.QueryRow("SELECT COUNT(*) FROM bans WHERE ip = $1", ip).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}
