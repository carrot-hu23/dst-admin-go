package DstAcfFileUtils

import (
	"dst-admin-go/utils/fileUtils"
	"log"
	"strconv"
	"strings"
)

type WorkshopItem struct {
	TimeUpdated int64
	Manifest    string
	Ugchandle   string
}

func ParseACFFile(filePath string) map[string]WorkshopItem {

	lines, err := fileUtils.ReadLnFile(filePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	parsingWorkshopItemsInstalled := false
	workshopItems := make(map[string]WorkshopItem)
	var currentItemID string
	var currentItem WorkshopItem
	for _, line := range lines {
		// log.Println(line)
		if strings.Contains(line, "WorkshopItemsInstalled") {
			parsingWorkshopItemsInstalled = true
			continue
		}

		if strings.Contains(line, "{") && parsingWorkshopItemsInstalled {
			continue
		}

		if strings.Contains(line, "}") {
			continue
		}

		if parsingWorkshopItemsInstalled {
			replace := strings.Replace(line, "\t\t", "", -1)
			replace = strings.Replace(replace, "\"", "", -1)
			if _, err := strconv.Atoi(replace); err == nil {
				// This line contains the Workshop Item ID
				currentItemID = line
			} else {
				// This line contains the Workshop Item details
				fields := strings.Fields(line)
				if len(fields) == 2 {
					key := strings.Replace(fields[0], "\"", "", -1)
					value := strings.Replace(fields[1], "\"", "", -1)
					// Remove double quotes from keys
					key = strings.ReplaceAll(key, "\"", "")
					switch key {
					case "timeupdated":
						currentItem.TimeUpdated, _ = strconv.ParseInt(value, 10, 64)
					case "manifest":
						currentItem.Manifest = strings.ReplaceAll(value, "\"", "")
					case "ugchandle":
						currentItem.Ugchandle = strings.ReplaceAll(value, "\"", "")
					}
				}
			}

			if currentItemID != "" && currentItem.TimeUpdated != 0 {
				workshopItems[currentItemID] = currentItem
				currentItemID = ""
				currentItem = WorkshopItem{}
			}
		}
	}

	return workshopItems
}
