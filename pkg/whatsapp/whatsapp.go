package whatsapp

import (
	"context"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	_ "modernc.org/sqlite"
)

var client *whatsmeow.Client

// Init initializes the WhatsApp client.
func Init(ctx context.Context, dbName *string) error {
	// Setup logging
	dbLog := waLog.Stdout("Database", "INFO", true)

	// Setup database for session storage
	var dbPath string
	if dbName == nil {
		dbPath = "file:whatsapp.db?_pragma=foreign_keys(1)"
	} else {
		dbPath = fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", *dbName)
	}

	// Changed from "sqlite3" to "sqlite" for modernc.org/sqlite driver
	container, err := sqlstore.New(ctx, "sqlite", dbPath, dbLog)
	if err != nil {
		return err
	}

	// Get first device (or create new one)
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return err
	}

	// Create WhatsApp client
	clientLog := waLog.Stdout("Client", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	// Register event handler
	client.AddEventHandler(eventHandler)

	// Connect to WhatsApp
	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(ctx)
		err = client.Connect()
		if err != nil {
			return err
		}

		// Print QR code for scanning
		fmt.Println("\n=== WhatsApp QR Code ===")
		fmt.Println("Please scan this QR code with WhatsApp on your phone:")
		fmt.Println("Settings → Linked Devices → Link a Device")

		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("\nWaiting for QR code scan...")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else if evt.Event == "success" {
				fmt.Println("✓ Successfully logged in!")
			} else {
				fmt.Printf("Login event: %s\n", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			return err
		}
		fmt.Println("✓ WhatsApp client connected (existing session)")
	}

	return nil
}

// Disconnect closes the WhatsApp connection
func Disconnect() {
	if client != nil {
		client.Disconnect()
		fmt.Println("WhatsApp client disconnected")
	}
}

// Event handler for incoming messages and events
func eventHandler(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Printf("Received message from %s: %s\n", v.Info.Sender, v.Message.GetConversation())

		// Example: Auto-reply to messages
		// SendMessage(context.Background(), v.Info.Sender.User, "Thanks for your message!")

	case *events.Receipt:
		if v.Type == types.ReceiptTypeRead || v.Type == types.ReceiptTypeReadSelf {
			fmt.Printf("%s read the message\n", v.Sender)
		} else if v.Type == types.ReceiptTypeDelivered {
			fmt.Println("Message delivered to", v.Sender)
		}

	case *events.Presence:
		if v.Unavailable {
			fmt.Printf("%s is now offline\n", v.From)
		} else {
			fmt.Printf("%s is now online\n", v.From)
		}

	case *events.LoggedOut:
		fmt.Println("⚠️  Logged out from WhatsApp!")
	}
}

// IsConnected checks if the client is connected
func IsConnected() bool {
	return client != nil && client.IsConnected()
}

// GetClient returns the WhatsApp client (for advanced usage)
func GetClient() *whatsmeow.Client {
	return client
}

// SendMessage sends a text message
func SendMessage(ctx context.Context, phoneNumber, message string) error {
	if client == nil {
		return fmt.Errorf("whatsapp client not initialized")
	}

	// Format: country code + phone number (without +)
	// Example: "254712345678" for Kenya
	jid := types.NewJID(phoneNumber, types.DefaultUserServer)

	msg := &waE2E.Message{
		Conversation: proto.String(message),
	}

	resp, err := client.SendMessage(ctx, jid, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	fmt.Printf("Message sent! ID: %s, Timestamp: %v\n", resp.ID, resp.Timestamp)
	return nil
}

// SendImage sends an image message
func SendImage(ctx context.Context, phoneNumber, imagePath, caption string) error {
	if client == nil {
		return fmt.Errorf("whatsapp client not initialized")
	}

	jid := types.NewJID(phoneNumber, types.DefaultUserServer)

	// Read image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to read image: %v", err)
	}

	// Upload image
	uploaded, err := client.Upload(ctx, imageData, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("failed to upload image: %v", err)
	}

	// Create image message
	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String("image/jpeg"),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(imageData))),
			Caption:       proto.String(caption),
		},
	}

	resp, err := client.SendMessage(ctx, jid, msg)
	if err != nil {
		return fmt.Errorf("failed to send image: %v", err)
	}

	fmt.Printf("Image sent! ID: %s\n", resp.ID)
	return nil
}

// SendDocument sends a document/file
func SendDocument(ctx context.Context, phoneNumber, filePath, fileName string) error {
	if client == nil {
		return fmt.Errorf("whatsapp client not initialized")
	}

	jid := types.NewJID(phoneNumber, types.DefaultUserServer)

	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Upload document
	uploaded, err := client.Upload(ctx, fileData, whatsmeow.MediaDocument)
	if err != nil {
		return fmt.Errorf("failed to upload document: %v", err)
	}

	// Create document message
	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String("application/pdf"),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(fileData))),
			FileName:      proto.String(fileName),
		},
	}

	resp, err := client.SendMessage(ctx, jid, msg)
	if err != nil {
		return fmt.Errorf("failed to send document: %v", err)
	}

	fmt.Printf("Document sent! ID: %s\n", resp.ID)
	return nil
}

// GetUserInfo gets user information
func GetUserInfo(ctx context.Context, phoneNumber string) error {
	if client == nil {
		return fmt.Errorf("whatsapp client not initialized")
	}

	jid := types.NewJID(phoneNumber, types.DefaultUserServer)

	info, err := client.GetUserInfo(ctx, []types.JID{jid})
	if err != nil {
		return err
	}

	for _, user := range info {
		fmt.Printf("User: %v\n", user)
	}
	return nil
}

// SendLocation sends a location message
func SendLocation(ctx context.Context, phoneNumber string, latitude, longitude float64) error {
	if client == nil {
		return fmt.Errorf("whatsapp client not initialized")
	}

	jid := types.NewJID(phoneNumber, types.DefaultUserServer)

	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(latitude),
			DegreesLongitude: proto.Float64(longitude),
		},
	}

	resp, err := client.SendMessage(ctx, jid, msg)
	if err != nil {
		return fmt.Errorf("failed to send location: %v", err)
	}

	fmt.Printf("Location sent! ID: %s\n", resp.ID)
	return nil
}
