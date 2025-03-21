package configs

import (
	// "crypto/rand"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() ConfigOTP {
	err := godotenv.Load() // Load .env file
	if err != nil {
		fmt.Println("Error loading .env file")
		createEnvFile()       // Create a new .env file
		err = godotenv.Load() // Reload the .env file after creation.
		if err != nil {
			fmt.Println("Error reloading .env file after creation")
			os.Exit(1)
		}
	}

	verifyToken := os.Getenv("VERIFY_TOKEN")
	secretKey := os.Getenv("SECRET_KEY")

	if verifyToken == "" || secretKey == "" {
		fmt.Println("Error: VERIFY_TOKEN or SECRET_KEY not found in .env")
		os.Exit(1)
	}

	return ConfigOTP{
		VerifyToken: verifyToken,
		SecretKey:   secretKey,
	}
}
func createEnvFile() {
	verifyToken, err := generateRandomString(32) // Độ dài tùy ý
	if err != nil {
		fmt.Println("Error generating VERIFY_TOKEN:", err)
		os.Exit(1)
	}

	secretKey, err := generateRandomString(64) // Độ dài tùy ý
	if err != nil {
		fmt.Println("Error generating SECRET_KEY:", err)
		os.Exit(1)
	}

	envContent := fmt.Sprintf("VERIFY_TOKEN=%s\nSECRET_KEY=%s\n", verifyToken, secretKey)

	err = os.WriteFile(".env", []byte(envContent), 0600)
	if err != nil {
		fmt.Println("Error writing .env file:", err)
		os.Exit(1)
	}

	fmt.Println(".env file created successfully!")
	fmt.Printf("VERIFY_TOKEN: %s\nSECRET_KEY: %s\n", verifyToken, secretKey) // In ra để đối chiếu
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
