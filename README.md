# RAG Chatbot with Go & ChromaDB 
This is a chatbot designed to help USF students answer any course-registration related questions they have in natural language using RAG!

<img width="940" alt="截屏2024-12-15 下午10 54 58_副本" src="https://github.com/user-attachments/assets/4a6cc9b8-9fd0-4f28-a1f6-e9f3c4a271b3" />


## YouTube Demo: 
https://www.youtube.com/watch?v=juo-J2kKP6w

## Tech Stack
- Go (amikos-tech's ChromaDB-in-Go and sashabaranov's OpenAI-in-Go packages)
- ChromaDB

## How the RAG system works:
Then chatbot will highlight specific fields that are related to course registration (names of professors, courses or subjects) and send back those fields in the form of a tool call in an attempt to query the database with those specific metadatas

e.g: 
- User enters "What courses does Jack Williams teach?" and sends the message to OpenAI
- OpenAI receives our message related to course selection, and sends back a structured JSON output that includes fields like: "InstructorFullName: 'Jack Williams'"
- Our program takes the fields highlighted by AI, and queries ChromaDB with those parameters ("InstructorFullName: 'Jack Williams'")
- By vector similarity, our program finds courses taught by professor "Christopher Brooks" (because professor Chris Brooks does not exist but Christopher Brooks does), and returns the information related to those courses in a string back to the chatbot as a response to the chatbot's tool call, which completes the tool call process.
- After the chatbot receives a response from its tool call, it then filters through the information it is given to answer the user's original question of "What courses does Jack Williams teach?" e.g: "Foundations of AI"

## How is the ChromaDB vector database built?
The vector database used to query course information for USF courses is built by parsing CSV rows into embeddings and metadatas.

Because all of the key fields related to a course's informaiton is known (CourseName, InstructorFullName, Building, MeetTime...), we can apply a set of metadatas for each course inserted into the collection in addition to parsing the course information into embeddings for more accurate queries.

An algorithm iterates the CSV by row, extracting the information and inserts the information into the database by batches.

## How to turn "fuzzy" instructor or subject names into canonical names?
Sometimes, the user doesn't know or will misspell the exact names of the things they are looking for. For instance, a professor may be well known as "Jack Williams", but his real name (stored in the university's catalogue) is "Jackson Williams".

In this situation, the query-by-metadata approach won't work, because there are no professors named "Jack Williams".

Therefore, before querying by metadatas, it is crucial to use **Vector Similary** to find the most likely canonical names of what the user is searching for. The program will account for this and search for canonical names before using them as literal metadatas for the final query.


