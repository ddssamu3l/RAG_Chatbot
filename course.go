package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Course struct {
    Subject                 string `json:"SUBJ"`
    CourseNumber            string `json:"CRSE NUM"`                        
    Section                 string `json:"SEC"`                      
    CRN                     string `json:"CRN"`                              
    ScheduleTypeCode        string `json:"Schedule Type Code"`               
    CampusCode              string `json:"Campus Code"`                    
    Title                   string `json:"Title Short Desc"`      
    InstructionMode         string `json:"Instruction Mode Desc"`           
    MeetingTypeCodes        string `json:"Meeting Type Codes"`              
    MeetDays                string `json:"Meet Days"`                       
    BeginTime               string `json:"Begin Time"`                        
    EndTime                 string `json:"End Time"`                        
    MeetStart               string `json:"Meet Start"`                       
    MeetEnd                 string `json:"Meet End"`                         
    Building                string `json:"BLDG"`                              
    Room                    string `json:"RM"`                                
    ActualEnrollment        string `json:"Actual Enrollment"`                 
    InstructorFirstName     string `json:"Primary Instructor First Name"`     
    InstructorLastName      string `json:"Primary Instructor Last Name"`      
    InstructorEmail         string `json:"Primary Instructor Email"`
    College                 string `json:"College"`                           
}

/*
courses := []*Course{}
if err := gocsv.Unmarshal(reader, &courses)
*/

func readCoursesFromCSV(filePath string) ([]Course, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open CSV file: %v", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.FieldsPerRecord = -1 

    // Read the header line
    headers, err := reader.Read()
    if err != nil {
        return nil, fmt.Errorf("failed to read CSV header: %v", err)
    }

    records, err := reader.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("failed to read CSV records: %v", err)
    }

    var courses []Course
    for _, record := range records {
		// err check in case of improperly formatted CSVs
        if len(record) < len(headers) {
            continue 
        }

        course := Course{
			Subject:             record[0],
			CourseNumber:        record[1],
			Section:             record[2],
			CRN:                 record[3],
			ScheduleTypeCode:    record[4],
			CampusCode:          record[5],
			Title:               record[6],
			InstructionMode:     record[7],
			MeetingTypeCodes:    record[8],
			MeetDays:            record[9],
			BeginTime:           record[10],
			EndTime:             record[11],
			MeetStart:           record[12],
			MeetEnd:             record[13],
			Building:            record[14],
			Room:                record[15],
			ActualEnrollment:    record[16],
			InstructorFirstName: record[17],
			InstructorLastName:  record[18],
			InstructorEmail:     record[19],
			College:             record[20],
		}
        courses = append(courses, course)
    }

    return courses, nil
}
