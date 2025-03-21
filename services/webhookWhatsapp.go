package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"otp/configs"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func HandleWhatsappWebhookRoutes(method, queryToken, challenge string, body []byte, config *configs.ConfigOTP, contractAddr common.Address, contractABI abi.ABI, client *ethclient.Client) (int, string) {
	switch method {
	case "GET":
		// Xử lý xác thực qua GET
		return HandleVerification(queryToken, challenge, config.VerifyToken)
	case "POST":
		// Xử lý webhook qua POST
		return HandleWebhook(body, config.SecretKey, contractAddr, contractABI, client)
	default:
		return 405, "Method not allowed"
	}
}

func SendMessageWhatsApp(to string, message string) error {
	// WhatsApp Business API URL
	url := fmt.Sprintf("https://graph.facebook.com/v13.0/%s/messages", configs.FromPhoneNumber)

	// Payload cho tin nhắn
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": message,
		},
	}

	// Convert payload thành JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error creating JSON payload: %v", err)
	}
	// Tạo request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	// Thêm headers
	req.Header.Set("Authorization", "Bearer "+configs.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Gửi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response from WhatsApp API:", string(body))
		return fmt.Errorf("error response from WhatsApp API: %s", body)
	}

	fmt.Println("Message sent successfully!")
	return nil
}

func HandleVerification(queryToken string, challenge string, verifyToken string) (int, string) {
	if queryToken == verifyToken {
		return 200, challenge
	} else {
		return 403, "Invalid verify token"
	}
}

func HandleWebhook(body []byte, secretKey string, contractAddr common.Address, contractABI abi.ABI, client *ethclient.Client) (int, string) {
	var payload configs.MessageWhatsapp
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 400, "Error parsing JSON"
	}

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			for _, message := range change.Value.Messages {
				otp := new(big.Int)
				if _, ok := otp.SetString(message.Text.Body, 10); !ok {
					fmt.Println("Invalid OTP format")
					return 400, "Invalid OTP format"
				}
				OTPVerified(contractAddr, contractABI, client, message.From, otp, 0, "WhatsApp")
			}
		}
	}
	return 200, "Webhook processed successfully"
}
