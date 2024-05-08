package comedy

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateSpanID() string {
	// Create a byte slice to store the random data
	randomBytes := make([]byte, 8)

	// Read random data from the crypto/rand package
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Error generating random bytes:", err)
		return ""
	}

	// Convert the random data to a hexadecimal string
	spanID := hex.EncodeToString(randomBytes)

	return spanID
}
