package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"otp/configs"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// sendTransaction creates, signs, and sends a transaction to the Ethereum network
func sendTransaction(client *ethclient.Client, contractAddress common.Address, privateKey *ecdsa.PrivateKey, data []byte, gasLimit uint64) (*types.Transaction, error) {
	publicKeyECDSA := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx, nil
}

func OTPVerified(contractAddress common.Address, contractABI abi.ABI, client *ethclient.Client, phoneNumber string, otp *big.Int, chatID int, TypeMethod string) {
	privateKey, err := crypto.HexToECDSA(configs.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	data, err := contractABI.Pack("validateOTP", otp, phoneNumber)
	if err != nil {
		log.Fatalf("Failed to pack function call: %v", err)
	}

	signedTx, err := sendTransaction(client, contractAddress, privateKey, data, 300000)
	if err != nil {
		log.Fatalf("Failed to send validateOTP transaction: %v", err)
	}

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Không lấy được số khối hiện tại: %v", err)
	}
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(blockNumber)),
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{contractABI.Events["OTPVerified"].ID}},
	}
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Không thể đăng ký lắng nghe logs: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("subscription error: %v", err)
		case vLog := <-logs:
			if vLog.TxHash == signedTx.Hash() {
				event := struct {
					PhoneNumber string
					Success     bool
					Message     string
					PublicKey   string
				}{}
				err := contractABI.UnpackIntoInterface(&event, "OTPVerified", vLog.Data)
				if err != nil {
					log.Fatalf("failed to unpack log: %v", err)
				}
				fmt.Printf("Event received - Phone: %s, Success: %v, Message: %s, PublicKey: %s\n",
					event.PhoneNumber, event.Success, event.Message, event.PublicKey)

				if !event.Success {
					return
				}

				txBytes, err := signedTx.MarshalBinary()
				if err != nil {
					log.Fatalf("Failed to marshal transaction to bytes: %v", err)
				}

				signature, err := signDataWithPrivateKey(privateKey, txBytes)
				if err != nil {
					log.Fatalf("Failed to sign data: %v", err)
				}

				encryptedData, err := encryptDataWithPublicKey(event.PublicKey, signature)
				if err != nil {
					log.Fatalf("Failed to encrypt data: %v", err)
				}
				fmt.Printf("Encrypted Data: %s\n", encryptedData)

				err = sendEncryptedDataToContract(client, contractABI, contractAddress, privateKey, encryptedData, event.PublicKey)
				if err != nil {
					log.Fatalf("Failed to send encrypted data to contract: %v", err)
				}

				sendMessage(TypeMethod, chatID, "OTP verification successful and data encrypted.")
				return
			}
		}
	}
}

// Sign data with the application's private key and return the signature
func signDataWithPrivateKey(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %v", err)
	}
	return signature, nil
}

// Hàm gửi dữ liệu đã mã hóa lên Smart Contract
func sendEncryptedDataToContract(client *ethclient.Client, contractABI abi.ABI, contractAddress common.Address, privateKey *ecdsa.PrivateKey, encryptedData string, publicKey string) error {
	data, err := contractABI.Pack("completeAuthentication", encryptedData, publicKey)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %v", err)
	}

	_, err = sendTransaction(client, contractAddress, privateKey, data, 300000)
	if err != nil {
		return fmt.Errorf("failed to send completeAuthentication transaction: %v", err)
	}

	return nil
}

func sendMessage(typeMethod string, chatID int, message string) {
	switch typeMethod {
	case "WhatsApp":
		SendMessageWhatsApp("84964928916", message)
	case "Telegram":
		SendMessageTelegram(chatID, message, configs.TelegramBotToken)
	default:
		// Xử lý khi không khớp với bất kỳ giá trị nào
	}
}

func encryptDataWithPublicKey(publicKeyStr string, data []byte) (string, error) {
	// modified := formatRSAPublicKey(publicKeyStr)
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block: invalid PEM format or empty")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("public key is not of RSA type")
	}
	keySize := rsaPublicKey.N.BitLen()
	if keySize < 2040 || keySize > 2056 {
		return "", fmt.Errorf("unsupported key size: expected ~2048 bits (2040–2056), got %d", keySize)
	}
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKey, data, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %v", err)
	}
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}
