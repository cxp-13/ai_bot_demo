package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Krognol/go-wolfram"
	"github.com/joho/godotenv"
	"github.com/slack-io/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
	"log"
	"os"
)

func main() {
	godotenv.Load(".env")

	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_TOKEN")}
	witAiClient := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.AddCommand(&slacker.CommandDefinition{
		Command:     "query for bot - <message>",
		Description: "send any question to wolfram",
		Examples:    []string{"who is the president of india"},
		Handler: func(ctx *slacker.CommandContext) {
			//userProfile := ctx.Event().UserProfile
			//fmt.Println(userProfile)
			query := ctx.Request().Param("message")

			msg, _ := witAiClient.Parse(&witai.MessageRequest{
				Query: query,
			})
			indent, _ := json.MarshalIndent(msg, "", "  ")

			fmt.Println(string(indent))
			value := gjson.Get(string(indent), "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			fmt.Printf("the value is %s\n", value.String())

			answerQuery, err := wolframClient.GetSpokentAnswerQuery(value.String(), wolfram.Metric, 1000)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Printf("the answer is %s\n", answerQuery)
			ctx.Response().Reply(answerQuery)
		},
	})

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
