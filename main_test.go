package main

import (
    "context"
    "encoding/json"
    "os"
    "testing"

    chroma "github.com/amikos-tech/chroma-go"
    "github.com/amikos-tech/chroma-go/pkg/embeddings/openai"
)

func TestVectorQuery(t *testing.T) {
    tests := []struct {
        name         string
        InstructorFirstName string
		InstructorLastName string
        expectedCRNs []string // Expected CRNs to validate the results
    }{
        {
            name:         "CS 272 by Philip Peterson",
            InstructorFirstName:  "Philip",
			InstructorLastName: "Peterson",
            expectedCRNs: []string{"40646"},
        },
	}
    ctx := context.Background()
    client, err := chroma.NewClient()
    if err != nil {
        t.Fatalf("Failed to create ChromaDB client: %v", err)
    }

    openaiEf, err := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
    if err != nil {
        t.Fatalf("Error creating OpenAI embedding function: %v", err)
    }

    collectionName := "usf-courses"

    collection, err := client.GetCollection(ctx, collectionName, openaiEf)
    if err != nil {
        t.Fatalf("Failed to get collection: %v", err)
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
			whereFilter := map[string]interface{}{
				"$or": []map[string]interface{}{
					{"InstructorFirstName": test.InstructorFirstName},
					{"InstructorLastName":  test.InstructorLastName},
				},
			}

            results, err := collection.Query(
                ctx,
                []string{test.InstructorFirstName, test.InstructorLastName},
                10,
                whereFilter,  
                nil,
                nil,
            )
            if err != nil {
                t.Fatalf("Error querying collection: %v", err)
            }

            if len(results.Documents) == 0 || len(results.Documents[0]) == 0 {
                t.Fatalf("No documents found for instructor: %s %s", test.InstructorFirstName, test.InstructorLastName)
            }

            foundCRNs := make(map[string]bool)
            for _, doc := range results.Documents {
                for _, retrievedDocument := range doc {
                    var retrievedCourse Course
                    err = json.Unmarshal([]byte(retrievedDocument), &retrievedCourse)
                    if err != nil {
                        t.Fatalf("Error unmarshaling retrieved document: %v", err)
                    }

                    foundCRNs[retrievedCourse.CRN] = true
                }
            }

            for _, expectedCRN := range test.expectedCRNs {
                if !foundCRNs[expectedCRN] {
                    t.Errorf("Expected CRN %s not found in the query results for '%s'", expectedCRN, test.name)
                }
            }
        })
    }
}