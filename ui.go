package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"github.com/sashabaranov/go-openai"
)

func StartUserInterface(db *Db){
	// openai client
	ctx := context.Background();
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	courseTool := MakeTool();
	emailTool := EmailTool();
	dialogue := InitializeDialogue()

	// scanner to take in command line inputs from user
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Search> ")
	for scanner.Scan(){
		question := scanner.Text()
		if question == "q"{
			return
		}
		
		dialogue = append(dialogue, openai.ChatCompletionMessage{
            Role:    openai.ChatMessageRoleUser,
            Content: question,
        })	

		resp, err := client.CreateChatCompletion(ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT4oMini,
				Messages: dialogue,
				Tools:    []openai.Tool{courseTool, emailTool},
			},
		)
		if err != nil {
			fmt.Printf("Completion error:%v len(choices):%v\n", err,
				len(resp.Choices))
			return	
		}
		if len(resp.Choices) == 0{
			fmt.Printf("No OpenAI response found. Skipping...\n")
			continue
		}

		msg := resp.Choices[0].Message
		dialogue = append(dialogue, msg)
		// if the AI called tools, handle the tools appropriately
		if len(msg.ToolCalls) != 0{
			for _, tool := range msg.ToolCalls {
				fmt.Printf("OpenAI called us back wanting to invoke our function '%v' with params '%v'\n",
					tool.Function.Name, tool.Function.Arguments)
				
				var content string
				// determine which tool is called and perform the appropriate action
				if tool.Function.Name == "email_instructor" {
					if err := emailProfessor(tool.Function.Arguments); err != nil {
						fmt.Printf("Error opening email: %v\n", err)
					}
					content = "an email draft has been opened successfully and the user has sent an email to the recipient."
				} else if tool.Function.Name == "get_relevant_courses" {
					whereFilter, err := BuildWhereFilterFromJSONString(db, tool.Function.Arguments)
					if err != nil {
						fmt.Printf("Error building WhereFilter: %v\n", err)
						continue
					} else {
						fmt.Printf("Trying to build WhereFilter with params: %v\nGot: %v\n\n", tool.Function.Arguments, whereFilter)
					}
		
					queryResults := queryDB(db, whereFilter)
		
					// add the chromaDB query results to our dialogue as a new chat message
					content = `If you believe you have enough information to answer the original user question with the information attached below, then answer it. Be sure to include all options to the user's question: ` + queryResults + "\n\nHowever, if you do not think you have enough information, then feel free to make another tool call."
				}
				// append the tool's response to our dialogue as a new chat message
				dialogue = append(dialogue, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    content,
					Name:       tool.Function.Name,
					ToolCallID: tool.ID,
				})
			}
			
			fmt.Printf("Sending OpenAI our function's response and requesting the reply to the original question...\n")
			resp, err = client.CreateChatCompletion(ctx,
				openai.ChatCompletionRequest{
					Model: openai.GPT4oMini,
					Messages: dialogue,
					Tools:    []openai.Tool{courseTool, emailTool},
				},
			)
			if err != nil || len(resp.Choices) != 1 {
				fmt.Printf("2nd completion error: err:%v len(choices):%v\n", err,
					len(resp.Choices))
				return
			}
			msg = resp.Choices[0].Message
		}
		
		
		// display OpenAI's response to the original question utilizing our function
		fmt.Printf("%v\n", msg.Content)
		fmt.Print("Search> ")
	}
}