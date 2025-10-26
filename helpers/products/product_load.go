package products

import (
	"database/sql"
)

type Product struct {
    ID       int
    Name     string
	Description string
    Weight   int
    Size     int
    Price    float64
    Quantity int
	Image string
	Image2 string
}

func LoadProducts(db *sql.DB) []Product{
	rows, _ := db.Query(`SELECT id, name, description, weight, size, price, quantity, image_url, image_url_2 FROM products`)
    var products []Product

    for rows.Next() {
        var p Product
        rows.Scan(&p.ID, &p.Name, &p.Description, &p.Weight, &p.Size, &p.Price, &p.Quantity, &p.Image, &p.Image2)
        products = append(products, p)
    }
	return products
}

