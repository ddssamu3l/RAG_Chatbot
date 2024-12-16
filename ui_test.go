package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestAIResponse(t *testing.T){
	tests := []struct{
		name string
		queryString string
		expected string
	}{
		{
			name: "TestPhilAndGreg",
			queryString: "What courses are Phil Peterson and Greg Benson teaching?",
			expected: `Here are the courses being taught by Phil Peterson and Greg Benson:

			### Phil Peterson
			1. **Course:** Software Development Lab
			   - **CRN:** 42344
			   - **Course Number:** 272L
			   - **Building:** MH
			   - **Room:** 122
			   - **Meet Days:** W
			   - **Begin Time:** 2:55 PM
			   - **End Time:** 4:25 PM
			   - **Instruction Mode:** In-Person
			
			2. **Course:** Software Development
			   - **CRN:** 40646
			   - **Course Number:** 272
			   - **Building:** LS
			   - **Room:** G12
			   - **Meet Days:** TR
			   - **Begin Time:** 2:40 PM
			   - **End Time:** 4:25 PM
			   - **Instruction Mode:** In-Person
			
			### Greg Benson
			1. **Course:** Computer Architecture
			   - **CRN:** 40649
			   - **Course Number:** 315
			   - **Building:** LS
			   - **Room:** 307
			   - **Meet Days:** TR
			   - **Begin Time:** 2:40 PM
			   - **End Time:** 4:25 PM
			   - **Instruction Mode:** In-Person
			
			2. **Course:** Laboratory
			   - **CRN:** 42346
			   - **Course Number:** 315L
			   - **Building:** LS
			   - **Room:** 307
			   - **Meet Days:** W
			   - **Begin Time:** 6:25 PM
			   - **End Time:** 7:55 PM
			   - **Instruction Mode:** In-Person
			
			3. **Course:** Laboratory
			   - **CRN:** 42345
			   - **Course Number:** 315L
			   - **Building:** LS
			   - **Room:** 307
			   - **Meet Days:** W
			   - **Begin Time:** 4:45 PM
			   - **End Time:** 6:15 PM
			   - **Instruction Mode:** In-Person
			
			4. **Course:** Computer Architecture
			   - **CRN:** 40649
			   - **Course Number:** 315
			   - **Building:** LS
			   - **Room:** 307
			   - **Meet Days:** TR
			   - **Begin Time:** 2:40 PM
			   - **End Time:** 4:25 PM
			   - **Instruction Mode:** In-Person`,
		},{
			name: "TestPHIL",
			queryString: "What philosophy courses are offered this semester?",
			expected: `You can take the following courses:
			
			1. **Great Philosophical Questions**
			   - CRN: 41168, SEC: 06
			   - Meeting Days: TR, Time: 1830 - 2015
			   - Instructor: Purushottama Bilimoria
			   - Room: KA 263
			
			2. **Great Philosophical Questions**
			   - CRN: 41167, SEC: 05
			   - Meeting Days: TR, Time: 1635 - 1820
			   - Instructor: Purushottama Bilimoria
			   - Room: KA 263
			
			3. **Philosophy of Biology**
			   - CRN: 41176, SEC: 01
			   - Meeting Days: MWF, Time: 1030 - 1135
			   - Instructor: Stephen Friesen
			   - Room: ED 102
			
			4. **Aesthetics**
			   - CRN: 41178, SEC: 02
			   - Meeting Days: MWF, Time: 1530 - 1635
			   - Instructor: Laurel Scotland-Stewart
			   - Room: KA 111
			
			5. **Great Philosophical Questions (Hybrid)**
			   - CRN: 41166, SEC: 04
			   - Meeting Days: F, Time: 1530 - 1635
			   - Instructor: Richie Kim
			   - Room: ONL
			
			6. **Mind, Freedom & Knowledge**
			   - CRN: 41196, SEC: 01
			   - Meeting Days: TR, Time: 1635 - 1820
			   - Instructor: Jennifer Fisher
			   - Room: ED 102
			
			7. **Environmental Ethics**
			   - CRN: 41192, SEC: 02
			   - Meeting Days: MWF, Time: 0915 - 1020
			   - Instructor: Stephen Friesen
			   - Room: ED 102
			
			8. **The Human Animal**
			   - CRN: 42008, SEC: 01
			   - Meeting Days: TR, Time: 1440 - 1625
			   - Instructor: Jennifer Fisher
			   - Room: ED 102
			
			9. **Existentialism**
			   - CRN: 41197, SEC: 01
			   - Meeting Days: TR, Time: 0955 - 1140
			   - Instructor: Brian Pines
			   - Room: ED 102
			
			10. **Ethics (Hybrid)**
				- CRN: 41188, SEC: 08
				- Meeting Days: MW, Time: 1415 - 1520
				- Instructor: Richie Kim
				- Room: KA 311
			
			11. **Philosophy of Science**
				- CRN: 41175, SEC: 02
				- Meeting Days: MW, Time: 1830 - 2015
				- Instructor: Krupa Patel
				- Room: KA 163
			
			12. **Philosophy of Religion**
				- CRN: 41172, SEC: 01
				- Meeting Days: MWF, Time: 1300 - 1405
				- Instructor: Deena Lin
				- Room: LM 244A
			
			13. **Logic**
				- CRN: 41199, SEC: 01
				- Meeting Days: MWF, Time: 1145 - 1250
				- Instructor: Nick Leonard
				- Room: ED 101
			
			14. **Ethics (Other Sections)**
				- Multiple sections and meeting times, taught by Greig Mulberry and Joshua Carboni. 
			
			15. **Topics in Contemporary Philosophy**
				- Multiple sections available, taught by various instructors.
			
			Please let me know if you need any more details on any specific course!`,
		},{
			name: "TestBio",
			queryString: "Where does Bioinformatics meet? Just say where the class meets.",
			expected: `The course titled "Bioinformatics" (CRSE NUM: 422) meets in room 311 of building KA, while the two sections of "Bioinformatics" (CRSE NUM: 640) meet in room 111 of building KA and room 136 of building HR. All classes are scheduled to meet on Mondays and Wednesdays (MW).`,
		},{
			name: "TestGuitar",
			queryString: "Can I learn guitar this semester?",
			expected: `Yes, you can learn guitar this semester by enrolling in the "Guitar and Bass Lessons" course (CRN: 41140 or 41141) taught by Christopher Ruscoe. The course is conducted in-person from August 20, 2024, to November 28, 2024.`,
		},{
			name: "TestMultiple",
			queryString: "I would like to take a Rhetoric course from Phil Choong. What can I take?",
			expected: `You can take the following Rhetoric courses taught by Philip Choong:

			1. **FYS: Podcasts: Eloquentia & Aud**
			   - CRN: 40215
			   - Course Number: 195
			   - Meeting Days: MWF
			   - Time: 2:15 PM - 3:20 PM
			   - Room: 352
			   - Actual Enrollment: 18
			
			2. **Public Speaking**
			   - CRN: 40146
			   - Course Number: 103
			   - Meeting Days: MWF
			   - Time: 10:30 AM - 11:35 AM
			   - Room: 346A
			   - Actual Enrollment: 22
			
			3. **Speaking Center Internship**
			   - CRN: 42533
			   - Course Number: 328
			   - Meeting Days: T
			   - Time: 4:35 PM - 6:25 PM
			   - Room: 345
			   - Actual Enrollment: 4
			
			Feel free to reach out to Philip Choong at pchoong@usfca.edu for more information!`,
		},
	}

	db, err := initializeDB()
    if err != nil {
        t.Fatalf("Error starting db: %v\n", err)
    }

    client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
    courseTool := MakeTool()
	emailTool := EmailTool();
    ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dialogue := InitializeDialogue()
			dialogue = append(dialogue, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: test.queryString,
			})	
	
			resp, err := client.CreateChatCompletion(ctx,
				openai.ChatCompletionRequest{
					Model: openai.GPT4oMini,
					Messages: dialogue,
					Tools:    []openai.Tool{courseTool, emailTool},
				},
			)
			if err != nil {
				t.Errorf("Completion error:%v len(choices):%v\n", err, len(resp.Choices))
			}
			if len(resp.Choices) == 0{
				t.Errorf("No OpenAI response found. Skipping...\n")
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
			
			aiFinalResponse := msg.Content
			
			if resp := compareAIResponseWithExpected(client, aiFinalResponse, test.expected); resp != "true"{
				t.Errorf(resp)
			}
		})
	}
}

func compareAIResponseWithExpected(client *openai.Client, aiResponse, expected string) string {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: `Your job is to compare the text below that contains information about some courses with the user-entered text. Check to see if some of key course information from the user-entered text (inside of the second message) is present in the first message. The naming of the data fields don't have to be exact. If the course information in the user-entered text is also present in the string immediately pasted below, then return the word "true". If not true, explain where the two texts differ. Here is the text you will compare the user-entered text to: ` + expected,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: aiResponse,
			},
		},
	}

	// send the comparison request
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil || len(resp.Choices) == 0 {
		return fmt.Sprintf("Comparison error: %v, choices length: %v", err, len(resp.Choices))
	}

	return resp.Choices[0].Message.Content
}
