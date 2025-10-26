package cart

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
)

type CartItem struct{
    ProductID int
    Image string
    Name string
    Quantity int
    Price float64
    Subtotal float64
    Stock int
}

func LoadCartItems(c *gin.Context, db *sql.DB) ([]CartItem,float64,error){
    
    cartID , err := GetCartID(c,db)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
		return []CartItem{}, 0, nil
	}
	    return nil, 0, err
    }

    rows, err := db.Query(`
    SELECT 
    cart_items.product_id,
    products.image_url,
    products.name,
    cart_items.quantity,
    products.price,
    products.quantity
    FROM 
    cart_items
    JOIN
    products ON cart_items.product_id = products.id
    WHERE 
    cart_items.cart_id = ?
    `,cartID)

    if err != nil {
        return []CartItem{},0,err
	}

    defer rows.Close()

    var items []CartItem
    var total float64

    for rows.Next() {
        var item CartItem
        err := rows.Scan(&item.ProductID, &item.Image, &item.Name, &item.Quantity, &item.Price, &item.Stock)
        if err != nil {
            continue
        }
        item.Subtotal = item.Price * float64(item.Quantity)
        total += item.Subtotal
        items = append(items, item)
    }

    return items,total,nil
}