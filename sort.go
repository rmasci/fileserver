package fileserver

import (
	"os"
	"sort"
)

func (dir *Directory) sortFSlist(fslist []string) []string {
	if dir.Sort == "" || (dir.Sort != "name" && dir.Sort != "size" && dir.Sort != "time") {
		dir.Sort = "time"
	}
	if dir.Direction == "" || (dir.Direction != "asc" && dir.Direction != "dec") {
		dir.Direction = "dec"
	}

	sort.Slice(fslist, func(i, j int) bool {
		switch dir.Sort {
		case "name":
			if dir.Direction == "asc" {
				return fslist[i] < fslist[j]
			}
			return fslist[i] > fslist[j]

		case "size":
			infoI, errI := os.Stat(fslist[i])
			infoJ, errJ := os.Stat(fslist[j])
			if errI != nil || errJ != nil {
				return false
			}
			sizeI := infoI.Size()
			sizeJ := infoJ.Size()
			if dir.Direction == "asc" {
				return sizeI < sizeJ
			}
			return sizeI > sizeJ
		case "time":
			fallthrough
		default:
			infoI, errI := os.Stat(fslist[i])
			infoJ, errJ := os.Stat(fslist[j])
			if errI != nil || errJ != nil {
				return false
			}
			timeI := infoI.ModTime()
			timeJ := infoJ.ModTime()
			if dir.Direction == "asc" {
				return timeI.Before(timeJ)
			}
			return timeI.After(timeJ)
		}
	})
	return fslist
}
