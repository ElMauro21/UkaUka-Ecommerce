package auth

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

func SendRecoveryEmail(c *gin.Context) error {
	
	emailRaw, emailExists := c.Get("reset_email")
	resetLink, linkExists := c.Get("reset_link")
	
	if !emailExists || !linkExists{
		return fmt.Errorf("context data missing: reset_email or reset_link not set")
	}

	email, ok := emailRaw.(string)
	if !ok {
		return fmt.Errorf("email in context is not a string")
	}
	
	password := os.Getenv("SMTP_PASSWORD")
  	if password == ""{
    	log.Fatal("SMTP_PASSWORD is not set")
  	}

	var(
		smtpHost = "smtp.gmail.com"
		smtpPort = "587"
		smtpUsername = "ukaukafb31@gmail.com"
	)
	
	from := smtpUsername
	to := []string{email}

	subject := "Subject: Recuperar contraseña\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	
	htmlBody := fmt.Sprintf(`
		<html>
  			<body>
    			<div style="
      				width: 100%%;
      				max-width: 600px;
      				height: 100px;
      				background-image: url('https://i.postimg.cc/xjXSMZYC/temp-Image-Gxql-NP.avif');
      				background-size: contain;
      	 			background-repeat: no-repeat;
      				background-position: left;
    			"></div>
    		<p>Hola,</p>
    		<p>Has solicitado recuperar tu contraseña. Haz clic en el siguiente botón para restablecerla:</p>
    		<a href="%s" style="display:inline-block;padding:10px 15px;background-color:rgb(210, 103, 51);color:white;text-decoration:none;border-radius:5px;">Restablecer contraseña</a>
   			<p>Si no solicitaste esto, puedes ignorar este mensaje.</p>
  			</body>
		</html>
	`, resetLink)

	message := []byte(subject + mime + htmlBody )

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}
	return nil
}

func SendRefundEmail(toEmail string, fullName string, productSummary string) error {
	password := os.Getenv("SMTP_PASSWORD")
	if password == "" {
		log.Fatal("SMTP_PASSWORD is not set")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUsername := "ukaukafb31@gmail.com"

	from := smtpUsername
	to := []string{toEmail}

	subject := "Subject: Reembolso de tu pedido en UkaUka\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	htmlBody := fmt.Sprintf(`
		<html>
  			<body>
				<div style="width:100%%;max-width:600px;padding:20px;font-family:Arial,sans-serif;">
					<h2 style="color:#d26733;">Hola %s,</h2>
					<p>Lamentamos informarte que uno o más productos de tu pedido están agotados:</p>
					<p><strong>%s</strong></p>
					<p>Como resultado, hemos iniciado un <strong>reembolso automático</strong> del valor correspondiente a tu método de pago.</p>
					<p>Gracias por tu comprensión. Si tienes dudas, contáctanos a ukaukafb31@gmail.com.</p>
					<p style="margin-top:30px;">— El equipo de <strong>UkaUka</strong></p>
				</div>
			</body>
		</html>
	`, fullName, productSummary)

	message := []byte(subject + mime + htmlBody)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}
	return nil
}

func SendRefundFailureEmail(email, fullName string) error {
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"
    smtpUsername := "ukaukafb31@gmail.com"
    password := os.Getenv("SMTP_PASSWORD")

    if password == "" {
        return fmt.Errorf("SMTP_PASSWORD is not set")
    }

    from := smtpUsername
    to := []string{email}

    subject := "Subject: Error al procesar tu reembolso\n"
    mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

    htmlBody := fmt.Sprintf(`
    <html>
      <body>
        <p>Hola %s,</p>
        <p>Intentamos procesar tu reembolso, pero ocurrió un error técnico.</p>
        <p>Estamos trabajando para resolverlo. Te notificaremos tan pronto como sea reembolsado correctamente.</p>
        <p>Gracias por tu paciencia.</p>
      </body>
    </html>`, fullName)

    message := []byte(subject + mime + htmlBody)

    auth := smtp.PlainAuth("", from, password, smtpHost)
    return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
}