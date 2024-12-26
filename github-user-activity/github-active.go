package main

import (
	"encoding/json" // JSON 데이터를 인코딩 및 디코딩하기 위한 패키지로, API 응답을 Go 구조체로 변환하거나 데이터를 JSON 형식으로 변환할 때 사용됩니다.
	"fmt"          // 문자열 포맷팅 및 출력과 관련된 기능을 제공하는 표준 패키지로, 터미널에 메시지를 출력할 때 사용됩니다.
	"io/ioutil"    // 파일이나 네트워크 응답에서 데이터를 읽을 때 사용하는 패키지로, API 응답 본문을 읽어오는 데 사용됩니다.
	"net/http"     // HTTP 요청 및 응답을 처리하기 위한 표준 패키지로, API에 GET 요청을 보내고 응답을 받는 데 사용됩니다.
	"os"           // 운영 체제와 상호작용하기 위한 패키지로, 명령줄 인자 처리 및 프로그램 종료에 사용됩니다.
)

// Activity 구조체는 하나의 GitHub 활동을 나타냅니다.
// Type: 활동의 유형(예: PushEvent, PullRequestEvent 등)
// Repo: 활동이 발생한 저장소 정보를 포함합니다.
// CreatedAt: 활동이 발생한 날짜 및 시간
// JSON 태그를 통해 GitHub API의 JSON 응답과 매핑됩니다.
type Activity struct {
	Type      string `json:"type"`
	Repo      struct {
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt string `json:"created_at"`
}

func main() {
	// 명령줄 인자로 사용자 이름이 전달되었는지 확인합니다.
	// os.Args는 명령줄 인자를 배열 형태로 저장하며, os.Args[0]은 실행 파일의 이름입니다.
	// 즉, 명령 구문이 짧게 입력될 시, Println 에 적힌 내용이 출력 됩니다.
	if len(os.Args) < 2 {
		fmt.Println("사용법: ./github-activity <GitHub 사용자 이름>")
		os.Exit(1) // 인자가 없을 경우 프로그램 종료
	}

	// 명령줄 인자로부터 사용자 이름을 가져옵니다.
	// os.Args[1]은 첫 번째 명령줄 인자를 의미합니다.
	username := os.Args[1]

	// GitHub API URL을 생성합니다.
	// username을 URL에 동적으로 삽입하여 해당 사용자의 활동을 가져옵니다.
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	// GitHub API에 GET 요청을 보냅니다.
	// http.Get은 URL에 대한 GET 요청을 수행하고 응답을 반환합니다.
	response, err := http.Get(url)
	if err != nil {
		// 요청 중 오류가 발생한 경우 메시지를 출력하고 종료합니다.
		fmt.Printf("요청 중 오류 발생: %v\n", err)
		os.Exit(1)
	}
	defer response.Body.Close() // 함수 종료 시 응답 본문을 닫아 리소스를 해제합니다.

	// 응답 상태 코드를 확인합니다.
	// HTTP 상태 코드가 200(OK)이 아닌 경우 오류로 간주합니다.
	if response.StatusCode != http.StatusOK {
		fmt.Printf("오류: 상태 코드 %d\n", response.StatusCode)
		os.Exit(1)
	}

	// 응답 본문을 읽습니다.
	// ioutil.ReadAll은 응답 본문 전체를 읽어 반환합니다.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("응답 본문 읽기 중 오류 발생: %v\n", err)
		os.Exit(1)
	}

	// JSON 응답을 Activity 구조체의 슬라이스로 변환합니다.
	// json.Unmarshal은 JSON 데이터를 Go 구조체로 디코딩합니다.
	var activities []Activity
	if err := json.Unmarshal(body, &activities); err != nil {
		fmt.Printf("JSON 파싱 중 오류 발생: %v\n", err)
		os.Exit(1)
	}

	// 활동 내용을 터미널에 출력합니다.
	// 활동의 날짜, 유형, 저장소 이름을 형식화하여 출력합니다.
	fmt.Printf("GitHub 사용자 '%s'의 최근 활동:\n\n", username)
	for i, activity := range activities {
		// 출력은 최근 10개의 활동으로 제한합니다.
		if i >= 10 {
			break
		}
		// CreatedAt: 활동 발생 날짜 및 시간
		// Type: 활동 유형
		// Repo.Name: 활동 저장소 이름
		fmt.Printf("[%s] %s in repository '%s'\n",
			activity.CreatedAt, activity.Type, activity.Repo.Name)
	}
}
