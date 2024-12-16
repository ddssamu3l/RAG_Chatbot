package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/pkg/embeddings/openai"
	"github.com/amikos-tech/chroma-go/types"
)

type Db struct {
	ctx                       context.Context
	client                    *chroma.Client
	coursesCollection         *chroma.Collection
	coursesCollectionName     string
	instructorsCollection     *chroma.Collection
	instructorsCollectionName string
	subjectsCollection        *chroma.Collection
	subjectsCollectionName    string
}

// handles the start process of the chromaDB db, getting/creating collections, and parsing data into the database
func Start(deleteFlag bool) (*Db, error) {
	db, err := initializeDB()
	if err != nil {
		log.Fatalf("Error creating database: %v\n", err)
		return nil, err
	}

	if deleteFlag {
		err := db.deleteCollections()
		if err != nil {
			log.Fatalf("Error deleting collections: %v", err)
			return nil, err
		}
		log.Printf("Successfully deleted collections")

		err = db.createCollections()
		if err != nil {
			return nil, fmt.Errorf("error creating collections: %w", err)
		}
		log.Printf("Successfully created collections")

		err = db.parseCSVIntoDatabase("Fall 2024 Class Schedule 08082024.csv")
		if err != nil {
			log.Fatalf("Error parsing CSV and/or inserting into database: %v\n", err)
			return nil, err
		}
	}

	return db, nil
}

// initialize the 'db' struct
func initializeDB() (*Db, error) {
	ct := context.Background()

	// Initialize ChromaDB client
	client, err := chroma.NewClient()
	if err != nil {
		log.Fatalf("Failed to create ChromaDB client: %v", err)
	}

	// Initialize OpenAI embedding function
	openaiEf, err := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatalf("Error creating OpenAI embedding function: %v", err)
	}

	// Define collection names
	coursesCollectionName := "usf-courses"
	instructorsCollectionName := "instructors"
	subjectsCollectionName := "subjects"

	// Get or create collections
	coursesCollection, err := client.GetCollection(ct, coursesCollectionName, openaiEf)
	if err != nil {
		if IsNotFoundError(err) {
			log.Printf("Courses collection '%s' does not exist, creating later.", coursesCollectionName)
			coursesCollection = nil
		} else {
			log.Fatalf("Failed to get courses collection: %v", err)
		}
	}

	instructorsCollection, err := client.GetCollection(ct, instructorsCollectionName, openaiEf)
	if err != nil {
		if IsNotFoundError(err) {
			log.Printf("Instructors collection '%s' does not exist, reating later.", instructorsCollectionName)
			instructorsCollection = nil
		} else {
			log.Fatalf("Failed to get instructors collection: %v", err)
		}
	}

	subjectsCollection, err := client.GetCollection(ct, subjectsCollectionName, openaiEf)
	if err != nil {
		if IsNotFoundError(err) {
			log.Printf("Subjects collection '%s' does not exist, creating later.", subjectsCollectionName)
			subjectsCollection = nil
		} else {
			log.Fatalf("Failed to get subjects collection: %v", err)
		}
	}

	db := Db{
		ctx:                       ct,
		client:                    client,
		coursesCollection:         coursesCollection,
		coursesCollectionName:     coursesCollectionName,
		instructorsCollection:     instructorsCollection,
		instructorsCollectionName: instructorsCollectionName,
		subjectsCollection:        subjectsCollection,
		subjectsCollectionName:    subjectsCollectionName,
	}

	return &db, nil
}

// deletes the three collections so that the database can remake these collections with new data
func (db *Db) deleteCollections() error {
	if db.client == nil {
		return fmt.Errorf("ChromaDB client is not initialized")
	}

	_, err := db.client.DeleteCollection(db.ctx, db.coursesCollectionName)
	if err != nil && !IsNotFoundError(err) {
		return fmt.Errorf("Error deleting courses collection: %v", err)
	}

	_, err = db.client.DeleteCollection(db.ctx, db.instructorsCollectionName)
	if err != nil && !IsNotFoundError(err) {
		return fmt.Errorf("Error deleting instructors collection: %v", err)
	}

	_, err = db.client.DeleteCollection(db.ctx, db.subjectsCollectionName)
	if err != nil && !IsNotFoundError(err) {
		return fmt.Errorf("Error deleting subjects collection: %v", err)
	}

	return nil
}

// create the courses, instructors, and subjects collections
func (db *Db) createCollections() error {
	if db.client == nil {
		return fmt.Errorf("ChromaDB client is not initialized")
	}

	openaiEf, err := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		return fmt.Errorf("error creating OpenAI embedding function: %w", err)
	}

	coursesCollection, err := db.client.CreateCollection(db.ctx, db.coursesCollectionName, nil, true, openaiEf, types.L2)
	if err != nil {
		return fmt.Errorf("failed to create courses collection '%s': %w", db.coursesCollectionName, err)
	}
	db.coursesCollection = coursesCollection

	instructorsCollection, err := db.client.CreateCollection(db.ctx, db.instructorsCollectionName, nil, true, openaiEf, types.L2)
	if err != nil {
		return fmt.Errorf("failed to create instructors collection '%s': %w", db.instructorsCollectionName, err)
	}
	db.instructorsCollection = instructorsCollection

	subjectsCollection, err := db.client.CreateCollection(db.ctx, db.subjectsCollectionName, nil, true, openaiEf, types.L2)
	if err != nil {
		return fmt.Errorf("failed to create subjects collection '%s': %w", db.subjectsCollectionName, err)
	}
	db.subjectsCollection = subjectsCollection

	return nil
}


func (db *Db) parseCSVIntoDatabase(filePath string) error {
	courses, err := readCoursesFromCSV(filePath)
	if err != nil {
		log.Fatalf("Error reading courses from CSV: %v", err)
	}

	var (
		courseDocuments     []string
		courseMetadatas     []map[string]interface{}
		courseIDs           []string
		instructorDocuments []string
		instructorIDs       []string
		subjectDocuments    []string
		subjectIDs          []string
	)

	instructorSet := make(map[string]struct{})
	subjectSet := make(map[string]struct{})

	for _, course := range courses {
		// Process courses collection
		courseJSON, err := json.Marshal(course)
		if err != nil {
			log.Fatalf("Error marshaling course to JSON: %v", err)
		}

		document := string(courseJSON)
		// Create metadata for querying
		metadata := map[string]interface{}{
			"CRN":                    course.CRN,
			"Subject":                course.Subject,
			"CourseNumber":           course.CourseNumber,
			"Section":                course.Section,
			"TitleShortDesc":         course.Title,
			"PrimaryInstructorEmail": course.InstructorEmail,
			"College":                course.College,
			"MeetDays":               course.MeetDays,
			"BeginTime":              course.BeginTime,
			"EndTime":                course.EndTime,
			"Building":               course.Building,
			"Room":                   course.Room,
			"InstructorFirstName":    course.InstructorFirstName,
			"InstructorLastName":     course.InstructorLastName,
			"InstructorFullName":     course.InstructorFirstName + " " + course.InstructorLastName,
		}

		// Generate a unique ID for the course using CRN
		id := course.CRN

		courseDocuments = append(courseDocuments, document)
		courseMetadatas = append(courseMetadatas, metadata)
		courseIDs = append(courseIDs, id)

		// gather all instructors in the csv without repeating names
		instructorFullName := course.InstructorFirstName + " " + course.InstructorLastName
		if _, exists := instructorSet[instructorFullName]; !exists {
			instructorSet[instructorFullName] = struct{}{}
			instructorDocuments = append(instructorDocuments, instructorFullName)
			instructorIDs = append(instructorIDs, instructorFullName)
		}

		// gather all subjects in the csv without repeating
		subject := course.Title
		if _, exists := subjectSet[subject]; !exists {
			subjectSet[subject] = struct{}{}
			subjectDocuments = append(subjectDocuments, subject)
			subjectIDs = append(subjectIDs, subject)
		}
	}

	// Insert into courses collection
	batchSize := 500
	for i := 0; i < len(courseDocuments); i += batchSize {
		end := i + batchSize
		if end > len(courseDocuments) {
			end = len(courseDocuments)
		}

		_, err = db.coursesCollection.Add(
			db.ctx,
			nil,
			courseMetadatas[i:end],
			courseDocuments[i:end],
			courseIDs[i:end],
		)
		if err != nil {
			log.Fatalf("Error adding documents to courses collection: %v", err)
		}
	}
	fmt.Printf("Successfully added %d courses to the courses collection.\n", len(courseDocuments))

	// Insert into instructors collection
	for i := 0; i < len(instructorDocuments); i += batchSize {
		end := i + batchSize
		if end > len(instructorDocuments) {
			end = len(instructorDocuments)
		}

		_, err = db.instructorsCollection.Add(
			db.ctx,
			nil,
			nil,
			instructorDocuments[i:end],
			instructorIDs[i:end],
		)
		if err != nil {
			log.Fatalf("Error adding documents to instructors collection: %v", err)
		}
	}
	fmt.Printf("Successfully added %d instructors to the instructors collection.\n", len(instructorDocuments))

	// Insert into subjects collection
	for i := 0; i < len(subjectDocuments); i += batchSize {
		end := i + batchSize
		if end > len(subjectDocuments) {
			end = len(subjectDocuments)
		}

		_, err = db.subjectsCollection.Add(
			db.ctx,
			nil,
			nil,
			subjectDocuments[i:end],
			subjectIDs[i:end],
		)
		if err != nil {
			log.Fatalf("Error adding documents to subjects collection: %v", err)
		}
	}
	fmt.Printf("Successfully added %d subjects to the subjects collection.\n", len(subjectDocuments))

	return nil
}

func IsNotFoundError(err error)bool{
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "does not exist")
}
