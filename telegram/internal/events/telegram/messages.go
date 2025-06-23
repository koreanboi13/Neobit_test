package telegram

const msgHelp = `
/test - measure the connection speed to Telegram servers.
/help - show this help message.
`

const msgHello = "Hello! I can help you test your connection speed to Telegram. Use /test to begin."

const (
	msgUnknownCommand  = "Unknown command. Use /help to see available commands."
	msgSpeedTestStart  = "Starting speed test... This may take a moment."
	msgSpeedTestFailed = "Sorry, the speed test failed."
)
