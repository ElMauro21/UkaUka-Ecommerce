package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/ElMauro21/UkaUkafb/helpers/cart"
	"github.com/ElMauro21/UkaUkafb/helpers/flash"
	"github.com/ElMauro21/UkaUkafb/helpers/view"
	"github.com/gin-gonic/gin"
)

func HandleOpenCart(c *gin.Context, db *sql.DB){
	
	items,total,err := cart.LoadCartItems(c,db)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al cargar el carrito: " + err.Error())
		return
	}

	description := cart.BuildDescription(items)

	msg,msgType := flash.GetMessage(c)

	view.Render(c,http.StatusOK,"cart.html",gin.H{
		"title": "My cart",
		"Message": msg,
		"MessageType": msgType,
		"items": items,
		"Total": total,
		"Description": description,
	})
}

func HandleAddToCart(c *gin.Context, db *sql.DB){
	
	if err := cart.CreateCart(c,db); err != nil {
		c.String(http.StatusInternalServerError, "Error al crear carrito: " + err.Error())
		return
	}
	
	productID := c.PostForm("product-id")
	quantity := c.PostForm("quantity")

	prodID,err := strconv.Atoi(productID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Id de producto invalido.")
		return
	}

	qty,err := strconv.Atoi(quantity)
	if err != nil {
		c.String(http.StatusInternalServerError, "Cantidad de producto invalido.")
		return
	}

	var stock int
	err = db.QueryRow(`SELECT quantity FROM products WHERE id = ?`,prodID).Scan(&stock)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se ha podido verificar el stock del producto.")
		return
	}

	cartID,err := cart.GetCartID(c,db)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se ha encontrado carrito de compras.")
		return
	}

	var currentQuantity int
	err = db.QueryRow(`SELECT quantity FROM cart_items WHERE cart_id = ? AND product_id = ?`,cartID,prodID).Scan(&currentQuantity)
	
	if errors.Is(err, sql.ErrNoRows) {
		currentQuantity = 0
		
		if currentQuantity+qty > stock {
			flash.SetMessage(c,"No se pueden añadir más productos de los que hay en stock!","error")
			c.Redirect(http.StatusSeeOther,"/shop")
			return
		}
		
		_,err := db.Exec(`INSERT INTO cart_items (cart_id, product_id, quantity) VALUES (?, ?, ?)`, cartID,prodID,qty)
		if err != nil {
			c.String(http.StatusInternalServerError, "No se ha podido añadir producto.")
			return
		}
	}else if err != nil {
		c.String(http.StatusInternalServerError, "Error al buscar el producto.")
	}else {

		if currentQuantity+qty > stock {
			flash.SetMessage(c,"No se pueden añadir más productos de los que hay en stock!","error")
			c.Redirect(http.StatusSeeOther,"/shop")
			return
		}

		_,err := db.Exec(`UPDATE cart_items SET quantity = quantity + ? WHERE cart_id = ? AND product_id = ?`,qty,cartID,prodID)
		if err != nil {
			c.String(http.StatusInternalServerError, "No se ha podido actualizar producto.")
			return
		}
	}

	flash.SetMessage(c,"Producto añadido al carrito","success")
	c.Redirect(http.StatusSeeOther,"/shop")
}

func HandleIncreaseQuantityCart(c *gin.Context, db *sql.DB) {
	productID := c.PostForm("product-id")
	prodID, err := strconv.Atoi(productID)
	if err != nil {
		c.String(http.StatusInternalServerError, "ID de producto inválido.")
		return
	}

	cartID, err := cart.GetCartID(c, db)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo obtener el carrito.")
		return
	}

	var currentQty int
	err = db.QueryRow(`
		SELECT quantity FROM cart_items WHERE cart_id = ? AND product_id = ?
	`, cartID, prodID).Scan(&currentQty)
	if err != nil {
		c.String(http.StatusInternalServerError, "Producto no encontrado en el carrito.")
		return
	}

	var stock int
	err = db.QueryRow(`SELECT quantity FROM products WHERE id = ?`, prodID).Scan(&stock)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al verificar el stock.")
		return
	}

	if currentQty+1 <= stock {
		_, err = db.Exec(`
			UPDATE cart_items SET quantity = quantity + 1 WHERE cart_id = ? AND product_id = ?
		`, cartID, prodID)
		if err != nil {
			c.String(http.StatusInternalServerError, "No se pudo actualizar la cantidad.")
			return
		}
		currentQty++
	}

	var price float64
	var name, image string
	err = db.QueryRow(`
		SELECT p.price, p.name, p.image_url
		FROM products p
		WHERE p.id = ?
	`, prodID).Scan(&price, &name, &image)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al obtener datos del producto.")
		return
	}

	items, total, err := cart.LoadCartItems(c, db)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al actualizar carrito")
		return
	}

	updatedItem := cart.FindCartItemByID(items, prodID)
	if updatedItem == nil {
		c.String(http.StatusNotFound, "Producto no encontrado después de actualizar.")
		return
	}

	description := cart.BuildDescription(items)

	view.Render(c, http.StatusOK, "cart_item_with_total.html", gin.H{
		"Item":  updatedItem,
		"Total": total,
		"Description": description,
	})

}

func HandleDecreaseQuantityCart(c *gin.Context, db *sql.DB){
	productID := c.PostForm("product-id")
	prodID, err := strconv.Atoi(productID)
	if err != nil {
		c.String(http.StatusInternalServerError, "ID de producto inválido.")
		return
	}

	cartID, err := cart.GetCartID(c, db)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo obtener el carrito.")
		return
	}

	var currentQty int
	err = db.QueryRow(`
		SELECT quantity FROM cart_items WHERE cart_id = ? AND product_id = ?
	`, cartID, prodID).Scan(&currentQty)
	if err != nil {
		c.String(http.StatusInternalServerError, "Producto no encontrado en el carrito.")
		return
	}

	var stock int
	err = db.QueryRow(`SELECT quantity FROM products WHERE id = ?`, prodID).Scan(&stock)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al verificar el stock.")
		return
	}

	if currentQty > 1 {
		_, err = db.Exec(`
			UPDATE cart_items SET quantity = quantity - 1 WHERE cart_id = ? AND product_id = ?
		`, cartID, prodID)
		if err != nil {
			c.String(http.StatusInternalServerError, "No se pudo actualizar la cantidad.")
			return
		}
		currentQty--
	}else if currentQty == 1 {
		_, err = db.Exec(`
			DELETE FROM cart_items WHERE cart_id = ? AND product_id = ?
		`, cartID, prodID)
		if err != nil {
			c.String(http.StatusInternalServerError, "No se pudo eliminar producto del carrito.")
			return
		}
		c.Header("HX-Redirect","/cart")
		c.Status(http.StatusSeeOther)
		return
	}

	var price float64
	var name, image string
	err = db.QueryRow(`
		SELECT p.price, p.name, p.image_url
		FROM products p
		WHERE p.id = ?
	`, prodID).Scan(&price, &name, &image)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al obtener datos del producto.")
		return
	}

	items, total, err := cart.LoadCartItems(c, db)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al actualizar carrito")
		return
	}

	updatedItem := cart.FindCartItemByID(items, prodID)
	if updatedItem == nil {
		c.String(http.StatusNotFound, "Producto no encontrado después de actualizar.")
		return
	}

	description := cart.BuildDescription(items)	

	view.Render(c, http.StatusOK, "cart_item_with_total.html", gin.H{
		"Item":  updatedItem,
		"Total": total,
		"Description": description,
	})
}
