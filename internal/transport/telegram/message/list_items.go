package message

import (
	"sort"
	"strconv"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

func itemList(items []domain.Item) ([]string, []string) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	itemsNames, keys := make([]string, len(items)), make([]string, len(items))

	for i := range items {
		itemsNames[i] = items[i].Name

		if items[i].Amount > 1 {
			itemsNames[i] += " ( " + strconv.FormatUint(uint64(items[i].Amount), 10) + " )"
		}

		keys[i] = items[i].Name
	}

	return itemsNames, keys
}
