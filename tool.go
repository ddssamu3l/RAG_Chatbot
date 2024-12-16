package main

import (
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func MakeTool() openai.Tool {
	// describe the function & its inputs
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"CRN": {
				Type:        jsonschema.String,
				Description: "Course Reference Number",
			},
			"Subject": {
				Type:        jsonschema.String,
				Description: "Subject code, e.g. CS",
			},
			"CourseNumber": {
				Type:        jsonschema.String,
				Description: "Course number, e.g. 272",
			},
			"Section": {
				Type:        jsonschema.String,
				Description: "Section number",
			},
			"TitleShortDesc": {
				Type:        jsonschema.String,
				Description: "The subject of the course. e.g. Bioinformatics",
			},
			"PrimaryInstructorEmail": {
				Type:        jsonschema.String,
				Description: "Email of the primary instructor",
			},
			"College": {
				Type:        jsonschema.String,
				Description: "Name of the college",
			},
			"MeetDays": {
				Type:        jsonschema.String,
				Description: "Days of the week when the course meets",
			},
			"BeginTime": {
				Type:        jsonschema.String,
				Description: "Start time of the course",
			},
			"EndTime": {
				Type:        jsonschema.String,
				Description: "End time of the course",
			},
			"Building": {
				Type:        jsonschema.String,
				Description: "Building where the course is held, e.g Lo Schiavo or LS",
			},
			"Room": {
				Type:        jsonschema.String,
				Description: "Room number where the course is held, e.g G12",
			},
			"InstructorFirstName": {
				Type:        jsonschema.String,
				Description: "First name of the instructor",
			},
			"InstructorLastName": {
				Type:        jsonschema.String,
				Description: "Last name of the instructor",
			},
			"InstructorFullName": {
				Type:        jsonschema.String,
				Description: "Full name of the instructor",
			},
		},
	}
	f := openai.FunctionDefinition{
		Name:        "get_relevant_courses",
		Description: "Get the relevant metadata from the user's question reguarding course information.",
		Parameters:  params,
	}
	t := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
	return t
}

func EmailTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"email": {
				Type:        jsonschema.String,
				Description: "The email address of the instructor",
			},
		},
	}
	f := openai.FunctionDefinition{
		Name:        "email_instructor",
		Description: "Opens up an email draft to a given instructor",
		Parameters:  params,
	}
	t := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
	return t
}

func InitializeDialogue () []openai.ChatCompletionMessage{
	dialogue := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: `You are a an agent who takes a prompt from a user and extracts key information from the user's free text and calls the appropriate function tools provided to you with parameters you are able to extract from the user text. Here's a list that contains all of the fields that can be from the user's free text: {
				"CRN": "40646",
				"Subject": "CS",
				"CourseNumber": "272",
				"Section": "03",
				"TitleShortDesc": "skating",
				"PrimaryInstructorEmail": "optimus.prime@usf.edu",
				"College": "College of Computer Science",
				"MeetDays": "MWF",
				"BeginTime": "09:00 AM",
				"EndTime": "10:15 AM",
				"Building": "Engineering Hall",
				"Room": "101",
				"InstructorFirstName": "Optimus",
				"InstructorLastName": "Prime",
				"InstructorFullName": "Optimus Prime"
			  }
			  the 'TitleShortDesc' field should be treated as the subject of the sentence. For instance, the 'TitleShortDesc' for the sentence "I want to learn skateboarding" should be "skateboarding". Please only extract fields if you are confident that the field exists in the user's text.
			  You can use the provided function tools to fetch the required information: {
                "get_relevant_courses": and parameters are the extracted fields from the user's text.
				"email_instructor": takes in the "email" field from the user's text.
              }
			  If you feel like the user's question requires multiple different tool calls or repeated tool calls of the same tool, you can just send one tool call, wait for the too call response, and continuing making subsequent calls until you have all the required information.`,
		},
	}
	return dialogue
}