package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/ElMauro21/UkaUkafb/helpers/flash"
	"github.com/ElMauro21/UkaUkafb/helpers/products"
	"github.com/ElMauro21/UkaUkafb/helpers/view"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Transaction struct {
	ID            int
	ReferenceCode string
	TotalAmount   float64
	CreatedAt     string
	UserName      string
	UserSurname   string
	Shipping      ShippingInfo
	Products      []Product
	ProductDetails string
}

type ShippingInfo struct {
	FullName     string
	IDNumber     string
	Phone        string
	Email        string
	State        string
	City         string
	Neighborhood string
	Address      string
}

type Product struct {
	ID       int    
	Name     string  
	Quantity int     
	Price    float64 
}


func HandleOpenDashboard(c *gin.Context, db *sql.DB) {

	session := sessions.Default(c)
	isAdmin, ok := session.Get("isAdmin").(bool)

	if !ok || !isAdmin {
		flash.SetMessage(c, "Necesitas permisos de administrador", "error")
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return
	}

	rows, err := db.Query(`
		SELECT t.id, t.reference_code, t.total_amount, t.created_at, u.names, u.surnames
		FROM transactions t
		JOIN users u ON t.user_id = u.id
		WHERE t.status = 'completed' AND t.shipped = 0
	`)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al obtener las transacciones.")
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.ReferenceCode, &t.TotalAmount, &t.CreatedAt, &t.UserName, &t.UserSurname); err != nil {
			c.String(http.StatusInternalServerError, "Error al procesar las transacciones.")
			return
		}

		var shipping ShippingInfo
		err = db.QueryRow(`
			SELECT full_name, id_number, phone, email, state, city, neighborhood, address
			FROM shipping_info
			WHERE transaction_id = ?
		`, t.ID).Scan(&shipping.FullName, &shipping.IDNumber, &shipping.Phone, &shipping.Email, &shipping.State, &shipping.City, &shipping.Neighborhood, &shipping.Address)

		if err != nil {
			c.String(http.StatusInternalServerError, "Error al obtener la información de envío.")
			return
		}

		t.Shipping = shipping

		rowsProducts, err := db.Query(`
			SELECT p.name, ti.quantity, p.price
			FROM transaction_items ti
			JOIN products p ON ti.product_id = p.id
			WHERE ti.transaction_id = ?
		`, t.ID)

		if err != nil {
			c.String(http.StatusInternalServerError, "Error al obtener los productos de la transacción.")
			return
		}
		defer rowsProducts.Close()

		var products []Product
		for rowsProducts.Next() {
			var p Product
			if err := rowsProducts.Scan(&p.Name, &p.Quantity, &p.Price); err != nil {
				c.String(http.StatusInternalServerError, "Error al procesar los productos.")
				return
			}
			products = append(products, p)
		}

		t.Products = products

		transactions = append(transactions, t)
	}

	products := products.LoadProducts(db)

	msg, msgType := flash.GetMessage(c)

	for i := range transactions {
    var productDetails []string
    for _, product := range transactions[i].Products {
        productDetails = append(productDetails, 
            product.Name + " (" + fmt.Sprintf("%d", product.Quantity) + ") $" + fmt.Sprintf("%.2f", product.Price))
    }
    transactions[i].ProductDetails = strings.Join(productDetails, ", ")
}

view.Render(c, http.StatusOK, "dashboard.html", gin.H{
	"title":       "Dashboard",
	"Message":     msg,
	"MessageType": msgType,
	"products":    products,  
	"transactions": transactions,  
})
}

func HandleAddProduct(c *gin.Context, db *sql.DB){

	name := c.PostForm("product-name")
	description := c.PostForm("product-description")
	weight := c.PostForm("product-weight")
	size := c.PostForm("product-size")
	price := c.PostForm("product-price")
	quantity := c.PostForm("product-quantity")
	image := c.PostForm("product-image")
	image2 := c.PostForm("product-image-two")

	if name == "" || description == "" || weight == "" || size == "" || price == "" || quantity == "" || image == "" || image2 == ""{
		view.RenderFlash(c,http.StatusOK,"Todos los campos son obligatorios","info")
		return
	}

	_, err := db.Exec(`INSERT INTO products 
	(name, description, weight, size, price, quantity, image_url, image_url_2) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, 
	name, description, weight, size, price, quantity, image, image2)
	if err != nil{
		c.String(http.StatusInternalServerError, "No se pudo crear el producto.")
		return
	}

	flash.SetMessage(c,"Producto creado correctamente","success")
	c.Header("HX-Redirect","/admin/dashboard")
	c.Status(http.StatusSeeOther)
}

func HandleDeleteProduct(c *gin.Context,db *sql.DB){

	productId := c.PostForm("product-id")

	if productId == "" {
		view.RenderFlash(c,http.StatusOK,"No hay producto para eliminar","info")
		return
	}

	name := c.PostForm("product-name")
	description := c.PostForm("product-description")
	weight := c.PostForm("product-weight")
	size := c.PostForm("product-size")
	price := c.PostForm("product-price")
	quantity := c.PostForm("product-quantity")
	image := c.PostForm("product-image")
	image2 := c.PostForm("product-image-two")
	
	if name == "" || description == "" || weight == "" || size == "" || price == "" || quantity == "" || image == "" || image2 == ""{
		view.RenderFlash(c,http.StatusOK,"Todos los campos son obligatorios","info")
		return
	}

	_, err := db.Exec(`DELETE FROM products WHERE id = ?`,productId)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo eliminar el producto.")
	}

	flash.SetMessage(c,"Producto eliminado correctamente","success")
	c.Header("HX-Redirect","/admin/dashboard")
	c.Status(http.StatusSeeOther)
}

func HandleUpdateProduct(c *gin.Context, db *sql.DB){

	productId := c.PostForm("product-id")
	if productId == "" {
		view.RenderFlash(c,http.StatusOK,"No hay producto para actualizar","info")
		return
	}

	name := c.PostForm("product-name")
	description := c.PostForm("product-description")
	weight := c.PostForm("product-weight")
	size := c.PostForm("product-size")
	price := c.PostForm("product-price")
	quantity := c.PostForm("product-quantity")
	image := c.PostForm("product-image")
	image2 := c.PostForm("product-image-two")

	if name == "" || description == "" || weight == "" || size == "" || price == "" || quantity == "" || image == "" || image2 == ""{
		view.RenderFlash(c,http.StatusOK,"Todos los campos son obligatorios","info")
		return
	}

	_, err := db.Exec(`UPDATE products
	SET name = ?, description = ?, weight = ?, size = ?, price = ?, quantity = ?, image_url = ?, image_url_2 = ? 
	WHERE id = ?
	`, name, description, weight, size, price, quantity, image, image2, productId)

	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo actualizar el producto.")
		return
	}

	flash.SetMessage(c,"Producto actualizado correctamente","success")
	c.Header("HX-Redirect","/admin/dashboard")
	c.Status(http.StatusSeeOther)
}