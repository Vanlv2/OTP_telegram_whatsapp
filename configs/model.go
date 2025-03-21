package configs

type MessageInfo struct {
	TypeOTP     string
	MessageBody string
}

// MessagePayload chứa cấu trúc dữ liệu của webhook message.
type MessageWhatsapp struct {
	// Cấu trúc dữ liệu cụ thể phụ thuộc vào API của WhatsApp Business
	// Ví dụ:
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				Messages []struct {
					From string `json:"from"`
					Text struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

// Telegram chúa cấu trúc dữ liệu từ Telegram
type MessageTelegram struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

// Config chứa các thông tin cấu hình cần thiết.
type ConfigOTP struct {
	VerifyToken string
	SecretKey   string // Thêm secret key để xác minh tính toàn vẹn của request
}
