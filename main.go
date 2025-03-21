package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"otp/configs"
	"otp/services"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("wss://bsc-testnet-rpc.publicnode.com")
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}
	services.SendMessageTelegram(6678028152, "test", configs.TelegramBotToken)

	contractAddr := common.HexToAddress(configs.ContractAddress)
	contractABI, err := abi.JSON(strings.NewReader(configs.ContractABI))
	if err != nil {
		log.Fatalf("Error parsing ABI: %v", err)
	}

	var wg sync.WaitGroup

	http.HandleFunc("/webhook/telegram", func(w http.ResponseWriter, r *http.Request) {
		// Đọc body trước khi phản hồi
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		// Phản hồi cho client
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Telegram webhook received! Processing asynchronously."))

		// Xử lý trong goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			services.ProcessTelegramWebhook(body, contractAddr, contractABI, client)
		}()
	})

	http.HandleFunc("/webhook/whatsapp", func(w http.ResponseWriter, r *http.Request) {
		configOTP := configs.LoadConfig()

		// Xử lý phần GET
		queryToken := r.URL.Query().Get("hub.verify_token")
		challenge := r.URL.Query().Get("hub.challenge")

		// Đọc body cho POST request
		var body []byte
		if r.Method == http.MethodPost {
			var err error
			body, err = ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				http.Error(w, "can't read body", http.StatusBadRequest)
				return
			}
			r.Body.Close()
		}

		statusCode, response := services.HandleWhatsappWebhookRoutes(r.Method, queryToken, challenge, body, &configOTP, contractAddr, contractABI, client)

		// Phản hồi lại client
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	})
	log.Println("HTTP server is running on port 8082 for both Telegram and WhatsApp webhooks...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	wg.Wait()
}
