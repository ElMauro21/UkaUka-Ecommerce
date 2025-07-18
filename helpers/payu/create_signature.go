package payu

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func GenerateSignature(apiKey, merchantID,  referenceCode, amount, currency string) string {
	signatureRaw := fmt.Sprintf("%s~%s~%s~%s~%s", apiKey,merchantID,referenceCode,amount,currency)
	hash := md5.Sum([]byte(signatureRaw))
	return hex.EncodeToString(hash[:])
}