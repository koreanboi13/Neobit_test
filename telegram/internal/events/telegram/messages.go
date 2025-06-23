package telegram

const msgHelp = `
/test <app_id> <app_hash> [proxy_addr] - measure the connection speed.
/help - show this help message.
`

const msgHello = "Hello! I can help you test your connection speed to Telegram. Use /test to begin."

const (
	msgUnknownCommand  = "Unknown command. Use /help to see available commands."
	msgInvalidTestCmd  = "Invalid command format. Use: /test <app_id> <app_hash> [proxy_addr]"
	msgSpeedTestStart  = "Starting speed test... This may take a moment."
	msgSpeedTestFailed = "Sorry, the speed test failed."
)
