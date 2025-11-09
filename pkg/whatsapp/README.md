# WhatsApp Integration Setup Guide

## Installation

### 1. Install Dependencies

```bash
go get go.mau.fi/whatsmeow
go get modernc.org/sqlite
go get google.golang.org/protobuf/proto
```

### 2. For QR Code Display in Terminal (Optional)

```bash
go get github.com/skip2/go-qrcode
# OR
go get github.com/mdp/qrterminal/v3
```

## Usage

### Basic Setup

1. **Run the application:**

```bash
   go run main.go
```

2. **Scan QR Code:**
   - A QR code will be printed in the terminal
   - Open WhatsApp on your phone
   - Go to Settings → Linked Devices → Link a Device
   - Scan the QR code

3. **Session Persistence:**
   - The session is saved in `whatsapp.db`
   - Next time you run the app, it will auto-connect without QR code

## Code Examples

### Send a Text Message

```go
// Phone number format: country code + number (no + sign)
// Example: "254712345678" for Kenya
err := sendMessage("254712345678", "Hello from Go!")
if err != nil {
    log.Fatal(err)
}
```

### Send an Image

```go
err := sendImage("254712345678", "./image.jpg", "Check this out!")
if err != nil {
    log.Fatal(err)
}
```

### Send a Document

```go
err := sendDocument("254712345678", "./report.pdf", "Monthly Report")
if err != nil {
    log.Fatal(err)
}
```

### Handle Incoming Messages

The `eventHandler` function already handles incoming messages. Customize it:

```go
case *events.Message:
    sender := v.Info.Sender.User
    message := v.Message.GetConversation()

    // Auto-reply example
    if message == "hello" {
        sendMessage(sender, "Hi! How can I help?")
    }
```

## Advanced Features

### Send to Group

```go
// Get group JID first (from incoming message or group list)
groupJID := types.NewJID("groupid", types.GroupServer)

msg := &waProto.Message{
    Conversation: proto.String("Hello group!"),
}

client.SendMessage(context.Background(), groupJID, msg)
```

### Get Profile Picture

```go
func getProfilePicture(phoneNumber string) (string, error) {
    jid := types.NewJID(phoneNumber, types.DefaultUserServer)
    pic, err := client.GetProfilePictureInfo(jid, nil)
    if err != nil {
        return "", err
    }
    return pic.URL, nil
}
```

### Check if Number is on WhatsApp

```go
func isOnWhatsApp(phoneNumber string) (bool, error) {
    jid := types.NewJID(phoneNumber, types.DefaultUserServer)
    resp, err := client.IsOnWhatsApp([]string{phoneNumber})
    if err != nil {
        return false, err
    }

    if info, ok := resp[jid]; ok {
        return info.IsIn, nil
    }
    return false, nil
}
```

### Send Location

```go
func sendLocation(phoneNumber string, latitude, longitude float64) error {
    jid := types.NewJID(phoneNumber, types.DefaultUserServer)

    msg := &waProto.Message{
        LocationMessage: &waProto.LocationMessage{
            DegreesLatitude:  proto.Float64(latitude),
            DegreesLongitude: proto.Float64(longitude),
        },
    }

    _, err := client.SendMessage(context.Background(), jid, msg)
    return err
}
```

## Folder Structure Integration

```bash
myproject/
├── cmd/
│   └── whatsapp-bot/
│       └── main.go              # Main application
├── internal/
│   ├── whatsapp/
│   │   ├── client.go            # WhatsApp client setup
│   │   ├── handlers.go          # Event handlers
│   │   ├── messages.go          # Message sending functions
│   │   └── media.go             # Media handling
│   └── config/
│       └── config.go            # Configuration
├── pkg/
│   └── utils/
│       └── logger.go            # Logging utilities
├── whatsapp.db                  # Session database (auto-created)
└── go.mod
```

## Important Notes

1. **Phone Number Format:** Always use format `countrycode+number` without `+`
or spaces
   - ✅ Correct: `254712345678`
   - ❌ Wrong: `+254 712 345 678` or `0712345678`

2. **Session Management:**
   - The `whatsapp.db` file stores your session
   - Keep it secure and don't commit to Git
   - Add to `.gitignore`

3. **Rate Limiting:**
   - WhatsApp has rate limits
   - Don't send too many messages too quickly
   - Use delays between bulk messages

4. **Terms of Service:**
   - Follow WhatsApp's Terms of Service
   - Don't spam users
   - Get consent before sending automated messages

5. **Multi-Device:**
   - This uses WhatsApp's multi-device protocol
   - Your phone doesn't need to be online
   - But your phone must have WhatsApp installed and active

## Troubleshooting

### QR Code Not Showing

Install qrterminal and modify the code:

```go
import "github.com/mdp/qrterminal/v3"

for evt := range qrChan {
    if evt.Event == "code" {
        qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
    }
}
```

### Connection Issues

- Check internet connection
- Delete `whatsapp.db` and re-authenticate
- Ensure WhatsApp is active on your phone

### Messages Not Sending

- Verify phone number format
- Check if number is on WhatsApp
- Ensure you're not rate-limited
