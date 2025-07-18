package main

import (
	"log"
	"os"

	"github.com/ElMauro21/UkaUkafb/database"
	"github.com/ElMauro21/UkaUkafb/handlers"
	"github.com/ElMauro21/UkaUkafb/helpers/admin"
	"github.com/ElMauro21/UkaUkafb/jobs"
	"github.com/ElMauro21/UkaUkafb/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
  err := godotenv.Load()
  if err !=nil {
    log.Println("No .env file found")
  }

  secret := os.Getenv("SESSION_SECRET")
  if secret == ""{
    log.Fatal("SESSION_SECRET is not set")
  }

  store := cookie.NewStore([]byte(secret))
  store.Options(sessions.Options{
	Path:     "/",        // Cookie is valid for all paths
	MaxAge:   3600,       // 1 hour
	HttpOnly: true,       // Not accessible from JS (security)
	Secure:   false,      // Use true to security
  })

  r := gin.Default()

  r.Use(sessions.Sessions("mysession",store))
  r.Use(middleware.InjectTemplateData())

  // Load the templates 
  r.LoadHTMLGlob("templates/*.html")
  // Serve any static file
  r.Static("/static","./static")

  db := database.InitDB("./uka.db")
  defer db.Close()

  jobs.JobTokenCleanup(db)
  jobs.JobGuestCartCleanup(db)

  admin.CreateAdminUser(db)

  // Home routes
  r.GET("/",handlers.HandleOpenHome)

  // Auth routes
  r.GET("/auth/login", handlers.HandleOpenLogin)
  r.POST("/auth/logout",handlers.HandleLogout)
  r.POST("/auth/login",func(c *gin.Context){
    handlers.HandleLogin(c,db)
  })
  r.POST("/auth/register",func(c *gin.Context){
    handlers.HandleRegister(c,db)
  })
  r.POST("/auth/recover/initiate",func(c *gin.Context){
    handlers.HandleCreateRecoveryLink(c,db)
  })
  r.GET("/auth/recover/reset", handlers.HandleShowResetForm)
  r.POST("/auth/recover/reset", func(c *gin.Context) {
	  handlers.HandleResetPassword(c, db)
  })

  // User routes
  r.GET("/user/profile", func (c *gin.Context){
    handlers.HandleOpenProfile(c,db)    
  })
  r.POST("/user/update", func(c *gin.Context) {
	  handlers.HandleUpdateProfile(c, db)
  })
  r.POST("/user/password", func(c *gin.Context) {
	  handlers.HandleChangePassword(c, db)
  })
  r.POST("/user/eliminate",func(c *gin.Context) {
    handlers.HandleDeleteAccount(c,db)
  })

  // Dashboard
  r.GET("/admin/dashboard", func (c *gin.Context){
    handlers.HandleOpenDashboard(c,db)
  }) 
  r.POST("/admin/products/add", func (c *gin.Context){
    handlers.HandleAddProduct(c,db)
  })
  r.POST("/admin/products/delete", func (c *gin.Context){
    handlers.HandleDeleteProduct(c,db)
  })
  r.POST("/admin/products/update", func (c *gin.Context){
    handlers.HandleUpdateProduct(c,db)
  })

  // Shop
  r.GET("/shop", func (c *gin.Context){
    handlers.HandleOpenShop(c,db)
  })

  // Cart
  r.GET("/cart", func (c *gin.Context){
    handlers.HandleOpenCart(c,db)
  })
  r.POST("/cart/add",func (c *gin.Context){
    handlers.HandleAddToCart(c,db)
  })
  r.POST("/cart/increase",func (c *gin.Context){
    handlers.HandleIncreaseQuantityCart(c,db)
  })
   r.POST("/cart/decrease",func (c *gin.Context){
    handlers.HandleDecreaseQuantityCart(c,db)
  })

  // Checkout
  r.POST("/checkout", func (c *gin.Context){
    handlers.HandleOpenCheckout(c,db)
  })
  r.POST("/process-payment", func (c *gin.Context){
    handlers.HandleProcessPayment(c,db)
  })

  r.Run() 
}
