package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"otp/configs"
	"regexp"
	"strconv"
	"strings"

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

	// Lấy phoneNumber và otp từ tin nhắn
	phoneNumber, otpStr, err := getPhoneNumberAndOTP(update.Message.Text)
	if err != nil {
		log.Printf("Error extracting phone number and OTP: %v", err)
		return
	}

	// Chuyển otp từ string sang uint64
	otp, err := strconv.ParseUint(otpStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing OTP to uint64: %v", err)
		return
	}

	// Chuyển otp từ uint64 sang *big.Int
	otpBig := new(big.Int).SetUint64(otp)

	// Thực hiện xác thực OTP
	OTPVerified(contractAddr, contractABI, client, phoneNumber, otpBig, update.Message.Chat.ID, "Telegram")
}

func getPhoneNumberAndOTP(messger string) (string, string, error) {
	// Loại bỏ khoảng trống thừa và ký tự xuống hàng
	cleaned := strings.TrimSpace(messger)                // Xóa khoảng trống đầu/cuối
	cleaned = strings.Join(strings.Fields(cleaned), " ") // Chuẩn hóa khoảng trống giữa các từ

	// Tách chuỗi thành các phần
	parts := strings.Split(cleaned, " ")
	if len(parts) < 2 {
		log.Println("Error: String format is invalid!")
		return "", "", fmt.Errorf("invalid string format")
	}

	phoneNumber := parts[0]
	otp := parts[1]

	// Kiểm tra định dạng số điện thoại: Chỉ cần là chuỗi số, có thể bắt đầu bằng + hoặc 0
	phoneRegex := regexp.MustCompile(`^\+?\d+$`) // Chấp nhận số bắt đầu bằng + hoặc không, chỉ chứa chữ số
	if !phoneRegex.MatchString(phoneNumber) {
		log.Println("Error: Invalid phone number format!")
		return "", "", fmt.Errorf("invalid phone number format")
	}

	// Kiểm tra định dạng OTP: 6 chữ số
	otpRegex := regexp.MustCompile(`^\d{6}$`)
	if !otpRegex.MatchString(otp) {
		log.Println("Error: Invalid OTP format!")
		return "", "", fmt.Errorf("invalid OTP format")
	}

	log.Println("Success: Phone Number:", phoneNumber)
	log.Println("Success: OTP:", otp)
	return phoneNumber, otp, nil
}
