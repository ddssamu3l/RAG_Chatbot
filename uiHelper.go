package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func BuildWhereFilterFromJSONString(db *Db, jsonStr string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var orConditions []map[string]interface{}

	for key, value := range data {
		// skip instructor name and subject for now because we will look at it individually later
		if key == "InstructorFullName" || key == "TitleShortDesc" {
			continue
		}

		if strValue, ok := value.(string); ok {
			trimmed := strings.TrimSpace(strValue)
			if trimmed == "" {
				continue
			}
			orConditions = append(orConditions, map[string]interface{}{key: trimmed})
		} else {
			orConditions = append(orConditions, map[string]interface{}{key: value})
		}
	}

	// turn fuzzy instructor name into canonical name
	if instructorFullName, exists := data["InstructorFullName"]; exists {
		if strValue, ok := instructorFullName.(string); ok {
			trimmed := strings.TrimSpace(strValue)
			if trimmed != "" {
				queryResults, err := db.instructorsCollection.Query(
					db.ctx,
					[]string{trimmed},
					1,
					nil,
					nil,
					nil,
				)
				if err != nil {
					return nil, fmt.Errorf("error querying instructors collection: %w", err)
				}

				// check if the collection actually returned anything
				if len(queryResults.Documents) > 0 && len(queryResults.Documents[0]) > 0 {
					instructorFullName := queryResults.Documents[0][0]
					fmt.Println("Instructor canonical name: ", instructorFullName)

					orConditions = append(orConditions, map[string]interface{}{"InstructorFullName": instructorFullName})
				}
			}
		}
	}

	if titleShortDesc, exists := data["TitleShortDesc"]; exists {
		if strValue, ok := titleShortDesc.(string); ok {
			trimmed := strings.TrimSpace(strValue)
			if trimmed != "" {
				queryResults, err := db.subjectsCollection.Query(
					db.ctx,
					[]string{trimmed},
					1,
					nil,
					nil, 
					nil,  
				)
				if err != nil {
					return nil, fmt.Errorf("error querying subjects collection: %w", err)
				}

				if len(queryResults.Documents) > 0 && len(queryResults.Documents[0]) > 0 {
					subjectName := queryResults.Documents[0][0]
					fmt.Println("Canonical subject course name: ", subjectName)
					orConditions = append(orConditions, map[string]interface{}{"TitleShortDesc": subjectName})
				}
			}
		}
	}

	if len(orConditions) == 1{
		orConditions = append(orConditions, map[string]interface{}{"Section": "999"})
	}

	// Construct the final whereFilter map with the accumulated orConditions
	whereFilter := map[string]interface{}{
		"$or": orConditions,
	}

	return whereFilter, nil
}

func queryDB(db *Db, whereFilter map[string]interface{})string{
	// query the metadatas with WhereFilter
	results, err := db.coursesCollection.Query(
		db.ctx,
		[]string{"."},
		50,
		whereFilter,  
		nil,
		nil,
	)
	if err != nil {
		fmt.Printf("Error querying collection: %v\n", err)
		fmt.Print("Search> ")
		return ""
	}

	// turn the DB query results into a []string
	matchingCourses := []string{}
	for _, doc := range results.Documents {
		for _, retrievedDocument := range doc {
			var retrievedCourse map[string]interface{}
			err = json.Unmarshal([]byte(retrievedDocument), &retrievedCourse)
			if err != nil {
				fmt.Printf("Error unmarshaling retrieved document: %v\n", err)
				continue
			}

			delete(retrievedCourse, "vector")
			retrievedCourseJSON, _ := json.MarshalIndent(retrievedCourse, "", "  ")
			matchingCourses = append(matchingCourses, string(retrievedCourseJSON))
		}
	}
	return strings.Join(matchingCourses, "\n")
}

func emailProfessor(jsonStr string)error{
	// Construct the mailto URL
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	if _, exists := data["email"]; !exists {
		return fmt.Errorf("email parameter not found in JSON")
	}

	professorEmail := data["email"].(string)
	mailto := fmt.Sprintf("mailto:%s", professorEmail)

	// account for different operating systems
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", mailto)
	case "darwin":
		cmd = exec.Command("open", mailto)
	default:
		return fmt.Errorf("unsupported platform")
	}

	// Execute the command to open the mail client
	return cmd.Start()
}