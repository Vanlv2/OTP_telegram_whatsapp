# Quy trình xác thực sử dụng Smart Contract, Golang và WhatsApp

## 0. Chuẩn bị
* Tạo bot tại: https://developers.facebook.com/apps/creation/

## 1. User tạo lệnh yêu cầu xác thực lên SM (Smart Contract)

*   **Hành động:**
    *   Người dùng gửi một giao dịch (transaction) đến Smart Contract (SM) để yêu cầu xác thực.
    *   Trong giao dịch này, người dùng đính kèm một khóa công khai (public key, kèm với address,time) ngẫu nhiên mà do dapp sinh cặp key ngẫu nhiên.
*   **Chi tiết kỹ thuật:**
    *   Giao dịch này sẽ gọi một hàm (function) cụ thể trên Smart Contract, ví dụ: `requestAuthentication(publicKey,address,time,type)`.
    *   `publicKey` là dữ liệu khóa công khai được tạo ra bởi người dùng (ví dụ: sử dụng thư viện `crypto/rsa` trong Go hoặc các thư viện tương tự trong ngôn ngữ khác).
    *`time` thời gian gửi
    *`type` loại như tele, whatsapp, Email
    *   Smart Contract sẽ lưu trữ khóa công khai này(publickey, time,type) liên kết với địa chỉ của người dùng.

## 2. SM (Smart Contract) tiến hành tạo đoạn OTP ngẫu nhiên và chỉ định số phone chatbot cho user

*   **Hành động:**
    *   Smart Contract tạo ra một mã OTP (One-Time Password) ngẫu nhiên.
    *   Smart Contract chỉ định một số điện thoại chatbot theo type (đã được cấu hình trước) cho người dùng.
*   **Chi tiết kỹ thuật:**
    *   OTP có thể được tạo bằng cách sử dụng các hàm ngẫu nhiên của ngôn ngữ lập trình Smart Contract (ví dụ: `random()` trong Solidity).
    *   Số điện thoại chatbot có thể được lưu trữ trong biến của Smart Contract.
    *   Smart Contract sẽ phát ra một sự kiện (event) chứa OTP và số điện thoại chatbot,type, ví dụ: `AuthenticationRequested(otp, chatbotPhoneNumber,type)`.

## 3. User vào WhatsApp add số phone chatbot và gửi OTP ngẫu nhiên để xác thực

*   **Hành động:**
    *   Người dùng thêm số điện thoại chatbot vào danh bạ WhatsApp.
    *   Người dùng gửi mã OTP nhận được đến số điện thoại chatbot qua WhatsApp.

## 4. Webhook viết bằng Golang nhận tin nhắn OTP, tiến hành gọi lên SM để kiểm tra OTP ngẫu nhiên, nếu có thì lấy public key ngẫu nhiên để ký cho dữ liệu và mã hoá rồi gửi lại SM

*   **Hành động:**
    *   Webhook (được viết bằng Golang) nhận được tin nhắn WhatsApp chứa OTP từ người dùng.
    *   Webhook gọi một hàm trên Smart Contract để kiểm tra xem OTP có hợp lệ hay không.
    *   Nếu OTP hợp lệ, Webhook lấy khóa công khai ngẫu nhiên của người dùng từ Smart Contract.
    *   Webhook ký dữ liệu liên quan đến quá trình xác thực bằng khóa riêng tư của ứng dụng.
    *   Webhook mã hóa dữ liệu đã ký bằng khóa công khai ngẫu nhiên của người dùng.
    *   Webhook gửi dữ liệu đã ký và mã hóa trở lại Smart Contract.
*   **Chi tiết kỹ thuật:**
    *   Webhook lắng nghe các sự kiện từ WhatsApp Webhook API.
    *   Gọi hàm trên Smart Contract: `validateOTP(otp)`.
    *   Lấy khóa công khai: `getPublicKey(userAddress)`.
    *   Ký dữ liệu: sử dụng thư viện `crypto/rsa` để ký dữ liệu bằng khóa riêng tư của ứng dụng.
    *   Mã hóa dữ liệu: sử dụng thư viện `crypto/rsa` để mã hóa dữ liệu bằng khóa công khai của người dùng.
    *   Gửi dữ liệu đã ký và mã hóa trở lại Smart Contract: `completeAuthentication(signedAndEncryptedData)`.
 
## 5. SM lưu hash của dữ liệu, để user có thể gửi dữ liệu cho bên thứ 3. và bên thứ 3 valid trên SM xem hash có đúng ko

*   **Hành động:**
    *   Smart Contract nhận dữ liệu đã ký và mã hóa từ Webhook.
    *   Smart Contract tính toán hash của dữ liệu này và lưu trữ hash.
    *   Người dùng có thể gửi dữ liệu đã ký và mã hóa cho bên thứ ba.
    *   Bên thứ ba có thể gửi hash của dữ liệu đến Smart Contract để xác minh.
*   **Chi tiết kỹ thuật:**
    *   Tính toán hash: sử dụng hàm `keccak256()` trong Solidity.
    *   Lưu trữ hash: lưu trữ hash trong một biến của Smart Contract.
    *   Xác minh hash: hàm `verifyHash(dataHash)` so sánh hash được gửi bởi bên thứ ba với hash đã lưu trữ.

## Lưu ý quan trọng

*   **Xử lý lỗi:** Cần xử lý các trường hợp lỗi như OTP không hợp lệ, lỗi mã hóa/giải mã, lỗi kết nối với Smart Contract.




