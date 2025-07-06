package cart

import (
	"fmt"
	"strings"
)

func BuildDescription(items []CartItem) string {
	var parts []string
	for _, item := range items {
		parts = append(parts, fmt.Sprintf("%s x %d", item.Name, item.Quantity))
	}
	return strings.Join(parts, ", ")
}
