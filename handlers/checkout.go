package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ElMauro21/UkaUkafb/helpers/auth"
	"github.com/ElMauro21/UkaUkafb/helpers/cart"
	"github.com/ElMauro21/UkaUkafb/helpers/flash"
	"github.com/ElMauro21/UkaUkafb/helpers/payu"
	"github.com/ElMauro21/UkaUkafb/helpers/users"
	"github.com/ElMauro21/UkaUkafb/helpers/view"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HandleOpenCheckout(c *gin.Context, db *sql.DB){
	
	msg,msgType := flash.GetMessage(c)

	u := users.LoadUserInfo(c,db)

	total := c.PostForm("total")
	description := c.PostForm("description")

	view.Render(c,http.StatusOK,"checkout.html",gin.H{
		"title": "Checkout",
		"User":  u,
		"Message": msg,
		"MessageType": msgType,
		"Total":       total,
		"Description": description,
	})
}

func HandleProcessPayment(c *gin.Context, db *sql.DB){
	fullName := c.PostForm("name") + " " + c.PostForm("surname")
	idNumber := c.PostForm("id-number")
	phone := c.PostForm("phone")
	email := c.PostForm("mail")
	state := c.PostForm("state")
	city := c.PostForm("city")
	neighborhood := c.PostForm("neighborhood")
	address := c.PostForm("address")
	description := c.PostForm("description")

	totalStr := c.PostForm("total")
	totalAmount, err := strconv.ParseFloat(totalStr,64)
	if err != nil {
    	c.String(http.StatusBadRequest, "Monto total inválido")
    	return
	}

	var userID *int
	var sessionID *string
	emailSession := sessions.Default(c).Get("user")

	if emailSession != nil {
    	id, err := auth.GetUserId(c, db)
    	if err != nil {
        	c.String(http.StatusInternalServerError, "Error obteniendo usuario")
        	return
    	}
    	userID = &id
	}else {
		s := sessions.Default(c).Get("cart_session_id")
    	if sStr, ok := s.(string); ok {
        	sessionID = &sStr
    	}
	}

	refCode, err := payu.CreateTransaction(db,userID, sessionID,totalAmount)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al crear transacción")
		return
	}

	var transactionID int
	err = db.QueryRow(`SELECT id FROM transactions WHERE reference_code = ?`, refCode).Scan(&transactionID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al obtener ID de la transacción")
		return
	}

	err = payu.SaveShippingInfo(db,
	fullName,idNumber,phone,email,state,city,neighborhood,address,transactionID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al guardar información de envío")
		return
	}

	apiKey := os.Getenv("API_KEY")
  	if apiKey == ""{
    	log.Fatal("API_KEY is not set")
  	}

	merchantId := os.Getenv("MERCHANT_ID")
  	if merchantId == ""{
    	log.Fatal("MERCHANT_ID is not set")
  	}

	accountId := os.Getenv("ACCOUNT_ID")
	if accountId == ""{
		log.Fatal("ACCOUNT_ID is not set")
	}
	
	cartItems,_,err := cart.LoadCartItems(c,db)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al cargar el carrito")
		return
	}
	
	for _, item := range cartItems {
		var stock int
		err := db.QueryRow("SELECT quantity FROM products WHERE id = ?", item.ProductID).Scan(&stock)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error verificando stock")
			return
		}
		if item.Quantity > stock {
			flash.SetMessage(c, "El producto '"+item.Name+"' no tiene suficiente stock.", "error")
			c.Redirect(http.StatusSeeOther, "/cart")
			return
		}
	}


	err = payu.SaveTransactionItems(db,transactionID,cartItems)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al guardar los ítems de la transacción")
	return
}

	signature := payu.GenerateSignature(apiKey,merchantId,refCode,fmt.Sprintf("%.2f",totalAmount),"COP")

	payuEnv := os.Getenv("PAYU_ENV")
	testFlag := "1" // default to sandbox
	payuFormURL := "https://sandbox.checkout.payulatam.com/ppp-web-gateway-payu/"
	if payuEnv == "production" {
    	testFlag = "0"
		payuFormURL = "https://checkout.payulatam.com/ppp-web-gateway-payu/"
	}

	c.HTML(http.StatusOK, "payu_form.html", gin.H{
		"MerchantID":         merchantId,
  		"AccountID":          accountId,
  		"Description":        description,
  		"ReferenceCode":      refCode,
  		"Amount":             fmt.Sprintf("%.2f", totalAmount),
  		"Tax":                "0",
  		"TaxReturnBase":      "0",
  		"Currency":           "COP",
  		"Signature":          signature,
  		"Test":               testFlag,
		"PayUFormURL":        payuFormURL,
  		"BuyerEmail":         email,
  		"BuyerFullName":      fullName,
  		"BuyerDocumentType":  "CC",
  		"BuyerDocument":      idNumber,
  		"Telephone":          phone,
  		"ShippingAddress":    address,
  		"ShippingCity":       city,
  		"ShippingCountry":    "CO",
  		"ResponseURL":        "https://6d807902255a.ngrok-free.app/payu/response",
  		"ConfirmationURL":    "https://6d807902255a.ngrok-free.app/payu/confirmation",
	})
}

func HandlePayUConfirmation(c *gin.Context, db *sql.DB){
	referenceCode := c.PostForm("reference_sale")
	state := c.PostForm("state_pol") // should be "4" for approved

	if referenceCode == "" || state != "4" {
		c.String(http.StatusBadRequest, "Invalid or non-approved transaction")
		return
	}

	transactionIDPayU := c.PostForm("transaction_id")
	if transactionIDPayU == "" {
    	transactionIDPayU = c.PostForm("transactionId")
	}

	_, err := db.Exec(`UPDATE transactions SET payu_transaction_id = ? WHERE reference_code = ?`, transactionIDPayU, referenceCode)
	if err != nil {
    	log.Printf("Error saving PayU transaction ID: %v", err)
	}
	
	err = payu.ProcessSuccessfulTransaction(db, referenceCode,transactionIDPayU)
	if err != nil {
		log.Printf("Error processing transaction %s: %v", referenceCode, err)
		c.String(http.StatusInternalServerError, "Error")
		return
	}
	
	c.String(http.StatusOK, "OK")
}

func HandleOpenSuccess(c *gin.Context, db *sql.DB) {
    referenceCode := c.Query("referenceCode")

    if referenceCode == "" {
        view.Render(c, http.StatusBadRequest, "payment_failed.html", gin.H{
            "title":   "Transacción inválida",
            "Message": "No se pudo verificar la transacción.",
        })
        return
    }

    var status string
    err := db.QueryRow(`SELECT status FROM transactions WHERE reference_code = ?`, referenceCode).Scan(&status)
    if err != nil {
        view.Render(c, http.StatusInternalServerError, "payment_failed.html", gin.H{
            "title":   "Error",
            "Message": "No se pudo validar el estado de tu transacción.",
        })
        return
    }

    switch status {
    case "completed":

        session := sessions.Default(c)
        session.Delete("cart_session_id")
        session.Save()

        flash.SetMessage(c, "¡Gracias por tu compra! Tu pago fue aprobado.", "success")
        view.Render(c, http.StatusOK, "success.html", gin.H{
            "title": "Pago exitoso",
        })

    case "needs_refund":
        view.Render(c, http.StatusOK, "refund_pending.html", gin.H{
            "title":   "Reembolso en proceso",
            "Message": "Tu pago fue recibido, pero el producto ya no está disponible. " + "Hemos iniciado un reembolso automático. Te llegará un correo.",
        })

    case "refunded":
        view.Render(c, http.StatusOK, "refund_completed.html", gin.H{
            "title":   "Reembolso completado",
            "Message": "Tu pago fue reembolsado automáticamente. " + "Revisa tu correo para más información.",
        })

    default:
        view.Render(c, http.StatusOK, "payment_failed.html", gin.H{
            "title":   "Pago no exitoso",
            "Message": "Tu pago no fue aprobado o ocurrió un problema. Intenta de nuevo.",
        })
    }
}
