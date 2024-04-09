package internal

const ExampleTokenHmacSecret = "my_secret"

var (
	Port    int
	Channel string
)

type ChatMessage struct {
	Input string `json:"input"`
}
