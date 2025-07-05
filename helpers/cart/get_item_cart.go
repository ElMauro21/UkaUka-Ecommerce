package cart

func FindCartItemByID(items []CartItem, id int) *CartItem {
	for i := range items {
		if items[i].ProductID == id {
			return &items[i]
		}
	}
	return nil
}