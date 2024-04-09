package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/centrifuge-cli/internal"

	"github.com/centrifugal/centrifuge-go"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version global
var Message string

func init() {
}

// Init command
func Init(rootCmd *cobra.Command) {
	var command = &cobra.Command{
		Use:   "publish",
		Short: "publish to a channel",
		Run: func(cmd *cobra.Command, args []string) {
			endpoint := fmt.Sprintf("ws://localhost:%d/connection/websocket", internal.Port)
			client := centrifuge.NewJsonClient(
				endpoint,
				centrifuge.Config{
					// Sending token makes it work with Centrifugo JWT auth (with `secret` HMAC key).
					Token: connToken("49", 0),
				},
			)
			defer client.Close()

			client.OnConnecting(func(e centrifuge.ConnectingEvent) {
				log.Printf("Connecting - %d (%s)", e.Code, e.Reason)
			})
			client.OnConnected(func(e centrifuge.ConnectedEvent) {
				log.Printf("Connected with ID %s", e.ClientID)
			})
			client.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
				log.Printf("Disconnected: %d (%s)", e.Code, e.Reason)
			})

			client.OnError(func(e centrifuge.ErrorEvent) {
				log.Printf("Error: %s", e.Error.Error())
			})

			client.OnMessage(func(e centrifuge.MessageEvent) {
				log.Printf("Message from server: %s", string(e.Data))
			})
			client.OnSubscribed(func(e centrifuge.ServerSubscribedEvent) {
				log.Printf("Subscribed to server-side channel %s: (was recovering: %v, recovered: %v)", e.Channel, e.WasRecovering, e.Recovered)
			})
			client.OnSubscribing(func(e centrifuge.ServerSubscribingEvent) {
				log.Printf("Subscribing to server-side channel %s", e.Channel)
			})
			client.OnUnsubscribed(func(e centrifuge.ServerUnsubscribedEvent) {
				log.Printf("Unsubscribed from server-side channel %s", e.Channel)
			})

			client.OnPublication(func(e centrifuge.ServerPublicationEvent) {
				log.Printf("Publication from server-side channel %s: %s (offset %d)", e.Channel, e.Data, e.Offset)
			})
			client.OnJoin(func(e centrifuge.ServerJoinEvent) {
				log.Printf("Join to server-side channel %s: %s (%s)", e.Channel, e.User, e.Client)
			})
			client.OnLeave(func(e centrifuge.ServerLeaveEvent) {
				log.Printf("Leave from server-side channel %s: %s (%s)", e.Channel, e.User, e.Client)
			})

			err := client.Connect()
			if err != nil {
				log.Fatalln(err)
			}
			msg := &internal.ChatMessage{
				Input: Message,
			}
			data, _ := json.Marshal(msg)

			_, err = client.Publish(context.Background(), internal.Channel, data)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("published message")
		},
	}
	var flagName = "message"
	command.PersistentFlags().StringVar(&Message, flagName, "hello", fmt.Sprintf("[required] i.e. --%s=hello", flagName))
	viper.BindPFlag(flagName, rootCmd.PersistentFlags().Lookup(flagName))

	rootCmd.AddCommand(command)
}
func connToken(user string, exp int64) string {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(internal.ExampleTokenHmacSecret))
	if err != nil {
		panic(err)
	}
	fmt.Println("token: ", t)
	return t
}
