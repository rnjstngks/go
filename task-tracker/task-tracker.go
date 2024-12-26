package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// Task 구조체는 개별 작업에 대한 정보를 저장합니다
// ID: 작업의 고유 식별자
// Description: 작업 설명
// Status: 작업 상태 ("todo", "in-progress", "done")
// CreatedAt: 작업이 생성된 시간
// UpdatedAt: 작업이 마지막으로 수정된 시간
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// 작업 데이터를 저장할 파일 이름을 정의합니다
const tasksFile = "tasks.json"

// loadTasks 함수는 JSON 파일에서 작업 목록을 읽어옵니다
// 파일이 없으면 빈 작업 목록을 반환합니다
func loadTasks() ([]Task, error) {
	var tasks []Task
	file, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		// 파일이 존재하지 않을 경우 빈 목록을 반환
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return nil, err
	}

	// JSON 데이터를 작업 목록으로 변환
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// saveTasks 함수는 작업 목록을 JSON 파일에 저장합니다
func saveTasks(tasks []Task) error {
	// 작업 목록을 JSON 형식으로 변환
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	// 변환된 JSON 데이터를 파일에 저장
	return ioutil.WriteFile(tasksFile, data, 0644)
}

// addTask 함수는 새 작업을 생성하고 저장합니다
// description: 작업 설명
func addTask(description string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("작업 로드 중 오류 발생:", err)
		return
	}

	// 새로운 작업 생성 및 ID 부여
	task := Task{
		ID:          len(tasks) + 1,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks = append(tasks, task)

	// 업데이트된 작업 목록 저장
	err = saveTasks(tasks)
	if err != nil {
		fmt.Println("작업 저장 중 오류 발생:", err)
		return
	}

	fmt.Printf("작업이 성공적으로 추가되었습니다: %d\n", task.ID)
}

// updateTask 함수는 기존 작업의 설명을 수정합니다
// taskID: 수정할 작업의 ID
// newDescription: 새 작업 설명
func updateTask(taskID int, newDescription string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("작업 로드 중 오류 발생:", err)
		return
	}

	// ID로 작업을 찾아 설명 수정
	for i, task := range tasks {
		if task.ID == taskID {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now()
			err = saveTasks(tasks)
			if err != nil {
				fmt.Println("작업 저장 중 오류 발생:", err)
				return
			}
			fmt.Printf("작업 %d이(가) 성공적으로 수정되었습니다.\n", taskID)
			return
		}
	}

	fmt.Printf("작업 %d을(를) 찾을 수 없습니다.\n", taskID)
}

// deleteTask 함수는 작업을 ID로 삭제합니다
// taskID: 삭제할 작업의 ID
func deleteTask(taskID int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("작업 로드 중 오류 발생:", err)
		return
	}

	// ID에 해당하지 않는 작업만 필터링하여 새로운 목록 생성
	newTasks := []Task{}
	for _, task := range tasks {
		if task.ID != taskID {
			newTasks = append(newTasks, task)
		}
	}

	// 업데이트된 작업 목록 저장
	err = saveTasks(newTasks)
	if err != nil {
		fmt.Println("작업 저장 중 오류 발생:", err)
		return
	}

	fmt.Printf("작업 %d이(가) 성공적으로 삭제되었습니다.\n", taskID)
}

// changeStatus 함수는 작업 상태를 변경합니다
// taskID: 상태를 변경할 작업의 ID
// status: 새 작업 상태 (예: "todo", "in-progress", "done")
func changeStatus(taskID int, status string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("작업 로드 중 오류 발생:", err)
		return
	}

	// ID로 작업을 찾아 상태 변경
	for i, task := range tasks {
		if task.ID == taskID {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			err = saveTasks(tasks)
			if err != nil {
				fmt.Println("작업 저장 중 오류 발생:", err)
				return
			}
			fmt.Printf("작업 %d의 상태가 %s(으)로 변경되었습니다.\n", taskID, status)
			return
		}
	}

	fmt.Printf("작업 %d을(를) 찾을 수 없습니다.\n", taskID)
}

// listTasks 함수는 작업 목록을 출력합니다
// status: 필터링할 상태 (빈 문자열이면 모든 작업 출력)
func listTasks(status string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("작업 로드 중 오류 발생:", err)
		return
	}

	// 상태에 따라 작업 필터링 및 출력
	for _, task := range tasks {
		if status == "" || task.Status == status {
			fmt.Printf("[ID: %d] %s (상태: %s)\n", task.ID, task.Description, task.Status)
		}
	}
}

// main 함수는 명령줄 입력을 처리하고 적절한 함수를 호출합니다
func main() {
	if len(os.Args) < 2 {
		fmt.Println("사용법: task-cli <명령> [옵션]")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("사용법: task-cli add <작업 설명>")
			return
		}
		description := os.Args[2]
		addTask(description)

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("사용법: task-cli update <작업 ID> <새 작업 설명>")
			return
		}
		taskID, _ := strconv.Atoi(os.Args[2])
		newDescription := os.Args[3]
		updateTask(taskID, newDescription)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("사용법: task-cli delete <작업 ID>")
			return
		}
		taskID, _ := strconv.Atoi(os.Args[2])
		deleteTask(taskID)

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("사용법: task-cli mark-in-progress <작업 ID>")
			return
		}
		taskID, _ := strconv.Atoi(os.Args[2])
		changeStatus(taskID, "in-progress")

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("사용법: task-cli mark-done <작업 ID>")
			return
		}
		taskID, _ := strconv.Atoi(os.Args[2])
		changeStatus(taskID, "done")

	case "list":
		status := ""
		if len(os.Args) > 2 {
			status = os.Args[2]
		}
		listTasks(status)

	default:
		fmt.Println("알 수 없는 명령입니다. 사용 가능한 명령: add, update, delete, mark-in-progress, mark-done, list")
	}
}
