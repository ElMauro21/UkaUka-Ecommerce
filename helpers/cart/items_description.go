package cart

import (
	"fmt"
	"regexp"
	"strings"
)

func sanitize(input string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return reg.ReplaceAllString(input, "")
}

func BuildDescription(items []CartItem) string {
	var parts []string
	for _, item := range items {
		name := sanitize(item.Name)
		parts = append(parts, fmt.Sprintf("%s x %d", name, item.Quantity))
	}
	description := strings.Join(parts, ", ")
	if len(description) > 255 {
		return description[:252] + "..."
	}
	return description
}