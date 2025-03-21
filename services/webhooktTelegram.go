package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"otp/configs"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Hàm gửi tin nhắn trả lời
func SendMessageTelegram(chatID int, text string, token string) error {
	// API URL để gửi tin nhắn
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	// Cấu trúc request JSON
	message := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	// Chuyển đổi message thành JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}
	// Gửi POST request tới Telegram API
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// Hàm xử lý webhook

func ProcessTelegramWebhook(body []byte, contractAddr common.Address, contractABI abi.ABI, client *ethclient.Client) {
	fmt.Println("ProcessTelegramWebhook is called")

	// Bắt đầu xử lý logic chính
	var update configs.MessageTelegram
	err := json.Unmarshal(body, &update)
	if err != nil {
		log.Printf("Error parsing webhook JSON: %v", err)
		return
	}

	fmt.Println("Received message from Telegram:", update.Message.Text)

	otpBigInt, err := strconv.ParseUint(update.Message.Text, 10, 0)
	if err != nil {
		log.Printf("Error parsing OTP text to uint256: %v", err)
		return
	}

	otp := new(big.Int).SetUint64(otpBigInt)

	// Thực hiện xác thực OTP
	OTPVerified(contractAddr, contractABI, client, update.Message.From.Username, otp, update.Message.Chat.ID, "Telegram")
}
