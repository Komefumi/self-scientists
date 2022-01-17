package data

import (
	"fmt"
	"self-scientists/config"
)

type ThreadPayload struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	CategoryName string `json:"categoryName"`
}

type PostPayload struct {
	ThreadID uint `json:"thread_id"`
	// ReplyingToID string `json:"replying_to_id"`
	Content string `json:"content"`
}

func (thread ThreadPayload) validateForCreation() []string {
	errors := []string{}

	if len(thread.Title) == 0 || len(thread.Description) == 0 || len(thread.CategoryName) == 0 {
		errors = append(errors, "All three must be provided: Title, Description, and CategoryName")
	}

	if len(errors) > 0 {
		return errors
	}

	var categoriesFound int
	row := config.DB.QueryRow("SELECT COUNT(*) FROM categories WHERE identifier=$1", thread.CategoryName)
	err := row.Scan(&categoriesFound)
	if err != nil {
		panic("Unable to connect to DB")
	}
	if categoriesFound == 0 {
		errors = append(errors, "No category with provided name exists")
	}

	return errors
}

func (thread ThreadPayload) CreateThread(creator_id uint) (errors []string, internallyErrored bool) {
	internallyErrored = false
	errors = thread.validateForCreation()
	if len(errors) > 0 {
		return
	}
	sqlStatement := `
      INSERT INTO threads (title, description, category_identifier, creator_id)
      VALUES ($1, $2, $3, $4)
    `
	_, dbErr := config.DB.Exec(sqlStatement, thread.Title, thread.Description, thread.CategoryName, creator_id)
	if dbErr != nil {
		internallyErrored = true
	}
	return
}

func GetThreadById(thread_id uint) (threadData map[string]interface{}, internallyErrored bool) {
	threadData = nil
	internallyErrored = false
	sqlStatement := `
		SELECT * FROM threads WHERE id=$1
	`
	rows, dbErr := config.DB.Query(sqlStatement, thread_id)
	if dbErr != nil {
		internallyErrored = true
		return
	}

	mappedDataList, errMappingOverData := getMapListFromSQLRows(rows)

	if errMappingOverData != nil {
		internallyErrored = true
		return
	}

	if len(mappedDataList) > 0 {
		threadData = mappedDataList[0]
	} else {
		threadData = nil
	}

	return
}

func GetThreadsPageCount() (pageCount uint, internallyErrored bool) {
	pageCount = 0
	internallyErrored = false

	sqlStatement := `SELECT COUNT(*) FROM threads`

	row := config.DB.QueryRow(sqlStatement)

	errScan := row.Scan(&pageCount)

	if errScan != nil {
		internallyErrored = true
	}

	return
}

func GetThreadListByPage(pageNumber uint) (threadDataList []map[string]interface{}, internallyErrored bool) {
	threadDataList = []map[string]interface{}{}
	internallyErrored = false
	sqlStatementThreadSelect := `
		SELECT * FROM threads ORDER BY updated_at DESC, created_at DESC OFFSET $1 LIMIT $2
	`
	rowsThreadSelect, dbErrThreadSelect := config.DB.Query(sqlStatementThreadSelect, ((pageNumber - 1) * config.ThreadPaginationSize), config.ThreadPaginationSize)

	if dbErrThreadSelect != nil {
		internallyErrored = true
		return
	}

	mappedDataList, errMappingOverData := getMapListFromSQLRows(rowsThreadSelect)

	if errMappingOverData != nil {
		internallyErrored = true
		return
	}

	threadDataList = mappedDataList

	return
}

func (post PostPayload) validateNewPostData() (errors []string, internallyErrored bool) {
	errors = []string{}
	internallyErrored = false

	if len(post.Content) == 0 {
		errors = append(errors, "Content for post should not be empty")
	}
	retrievedThread, internallyErroredForThreadFetch := GetThreadById(post.ThreadID)
	if internallyErroredForThreadFetch {
		internallyErrored = true
		return
	}

	if retrievedThread == nil {
		errors = append(errors, "Thread specified by threadId not found")
	}

	return
}

func (post PostPayload) CreatePost(userId uint) (errors []string, internallyErrored bool) {
	internallyErrored = false
	errors = []string{}
	errors, internallyErroredAtValidation := post.validateNewPostData()

	if internallyErroredAtValidation {
		internallyErrored = true
		return
	}
	if len(errors) > 0 {
		return
	}

	sqlStatementPostCreation := `
		INSERT INTO posts (thread_id, author_id, content)
    VALUES ($1, $2, $3)
	`
	_, dbErr := config.DB.Exec(sqlStatementPostCreation, post.ThreadID, userId, post.Content)

	if dbErr != nil {
		internallyErrored = true
		return
	}

	return
}

func GetPostsListForThreadByPage(threadId uint, pageNumber uint) (threadData interface{}, postDataList []map[string]interface{}, internallyErrored bool) {
	threadData = nil
	postDataList = []map[string]interface{}{}
	internallyErrored = false

	retrievedThreadData, internallyErroredThreadRetrieval := GetThreadById(threadId)

	if internallyErroredThreadRetrieval {
		internallyErrored = true
		return
	}
	if retrievedThreadData == nil {
		return
	}

	threadData = retrievedThreadData
	sqlStatementPostListSelect := `
		SELECT * FROM posts WHERE thread_id=$1 ORDER BY updated_at DESC, created_at DESC OFFSET $2 LIMIT $3
	`
	rowsPostSelect, dbErrPostSelect := config.DB.Query(sqlStatementPostListSelect, threadId, ((pageNumber - 1) * config.PostPaginationSize), config.PostPaginationSize)

	if dbErrPostSelect != nil {
		fmt.Println(dbErrPostSelect)
		internallyErrored = true
		return
	}

	mappedPostDataList, errMappingOverPostDataList := getMapListFromSQLRows(rowsPostSelect)

	if errMappingOverPostDataList != nil {
		fmt.Println(errMappingOverPostDataList)
		internallyErrored = true
		return
	}

	postDataList = mappedPostDataList

	return
}
