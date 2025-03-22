package configs

// ngrok config add-authtoken 2uDHKB3y43fFC1G4jX6CZjuWM1A_4ipy24dhQvogqKcHmhg21
//https://developers.facebook.com/apps/1090363283107015/whatsapp-business/wa-dev-console/?business_id=211741613535919
//https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook?url=https://<YOUR_DOMAIN>/webhook/telegram

const FromPhoneNumber = "646787801840640" //id của bot whatsapp
const AccessToken = "EAAPfrfwqbMcBO2ZA0xphpRvTQg1TubDKQ984fmBf7FRVKSSyh04pDa8ube0N5wNpCl586r1ChUg07IrfQesZCPbpTp1FxU0AMiqQo5rC1xZBGNZCCl9CAALE4Wm20frzEyAC5Uz7sVT3PZCgcnv9rACT96HQRhRC0z3nr5tu6C0k3Vei5wJ9FvujdBX8S4Sa3a4SNlIticAxY6OZAeIuDYDaBoii0ZD"

const TelegramBotToken = "7746087401:AAFiG3CA9ZT_5u2prA0tRyCYrkAUy0qAcDs"

const PrivateKey = "1318c69b804769a5472a284f80853d4c54914fbf289be1bac2480f82886b43bf"
const ContractAddress = "0xd0E24776D6254E54EC3645105f6c07166a819f5e"
const ContractABI = `[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "user",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "otp",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "string",
				"name": "chatbotPhone",
				"type": "string"
			},
			{
				"indexed": false,
				"internalType": "enum AuthOTP.TypeMethod",
				"name": "typeMethod",
				"type": "uint8"
			}
		],
		"name": "AuthenticationRequested",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "string",
				"name": "phoneNumber",
				"type": "string"
			},
			{
				"indexed": false,
				"internalType": "bool",
				"name": "success",
				"type": "bool"
			},
			{
				"indexed": false,
				"internalType": "string",
				"name": "message",
				"type": "string"
			},
			{
				"indexed": false,
				"internalType": "string",
				"name": "publicKey",
				"type": "string"
			}
		],
		"name": "OTPVerified",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "_phoneNumber",
				"type": "string"
			},
			{
				"internalType": "enum AuthOTP.TypeMethod",
				"name": "_typeMethod",
				"type": "uint8"
			}
		],
		"name": "addBot",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "_data",
				"type": "string"
			},
			{
				"internalType": "string",
				"name": "_publicKey",
				"type": "string"
			}
		],
		"name": "completeAuthentication",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "detailBots",
		"outputs": [
			{
				"internalType": "string",
				"name": "phoneNumber",
				"type": "string"
			},
			{
				"internalType": "enum AuthOTP.TypeMethod",
				"name": "typeMethod",
				"type": "uint8"
			},
			{
				"internalType": "bool",
				"name": "busy",
				"type": "bool"
			},
			{
				"internalType": "uint256",
				"name": "timeOccupied",
				"type": "uint256"
			},
			{
				"internalType": "bool",
				"name": "status",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			}
		],
		"name": "publicKeyUsers",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "_userPhoneNumber",
				"type": "string"
			},
			{
				"internalType": "string",
				"name": "_publicKey",
				"type": "string"
			},
			{
				"internalType": "enum AuthOTP.TypeMethod",
				"name": "_typeMethod",
				"type": "uint8"
			}
		],
		"name": "requestAuthentication",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_botId",
				"type": "uint256"
			},
			{
				"internalType": "string",
				"name": "_phoneNumber",
				"type": "string"
			},
			{
				"internalType": "enum AuthOTP.TypeMethod",
				"name": "_typeMethod",
				"type": "uint8"
			},
			{
				"internalType": "bool",
				"name": "_status",
				"type": "bool"
			}
		],
		"name": "updateBot",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_otp",
				"type": "uint256"
			},
			{
				"internalType": "string",
				"name": "userPhoneNumber",
				"type": "string"
			}
		],
		"name": "validateOTP",
		"outputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "_publicKey",
				"type": "string"
			},
			{
				"internalType": "bytes32",
				"name": "_dataHash",
				"type": "bytes32"
			}
		],
		"name": "verifyHash",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`
