// Go 언어로 작성된 간단한 Expense Tracker CLI 프로그램

package main

import (
	"bufio"      // 표준 입력을 쉽게 읽을 수 있도록 지원하는 모듈
	"encoding/csv" // 데이터를 CSV 형식으로 읽고 쓰는 데 사용하는 모듈
	"fmt"         // 입출력 형식을 지원하는 표준 라이브러리
	"os"          // 파일 시스템 및 운영 체제와 상호 작용하기 위한 모듈
	"strconv"     // 문자열을 숫자로 변환하거나 숫자를 문자열로 변환하기 위한 모듈
	"strings"     // 문자열 조작을 위한 다양한 유틸리티를 제공하는 모듈
)

type Expense struct {
	Description string  // 지출 내역 설명
	Amount      float64 // 지출 금액
}

var expenses []Expense // 모든 지출 항목을 저장하는 슬라이스

// 메인 함수: 프로그램 실행의 진입점
func main() {
	reader := bufio.NewReader(os.Stdin) // 표준 입력에서 줄 단위로 데이터를 읽기 위해 bufio.Reader 생성
	for {
		fmt.Println("\n지출 관리 프로그램")
		fmt.Println("1. 지출 추가")
		fmt.Println("2. 지출 목록 보기")
		fmt.Println("3. CSV 파일로 내보내기")
		fmt.Println("4. 종료")
		fmt.Print("선택: ")

		input, _ := reader.ReadString('\n') // 사용자 입력을 읽음
		input = strings.TrimSpace(input)   // 입력값의 앞뒤 공백 제거

		switch input {
		case "1":
			addExpense(reader) // 지출 추가 함수 호출
		case "2":
			viewExpenses() // 지출 목록 보기 함수 호출
		case "3":
			exportToCSV() // CSV 파일 내보내기 함수 호출
		case "4":
			fmt.Println("프로그램을 종료합니다.")
			return
		default:
			fmt.Println("잘못된 입력입니다. 다시 시도하세요.")
		}
	}
}

// 지출 추가 함수
func addExpense(reader *bufio.Reader) {
	fmt.Print("지출 설명: ")
	desc, _ := reader.ReadString('\n') // 지출 설명 입력받음
	desc = strings.TrimSpace(desc)

	fmt.Print("지출 금액: ")
	amountStr, _ := reader.ReadString('\n') // 지출 금액 입력받음
	amountStr = strings.TrimSpace(amountStr)

	amount, err := strconv.ParseFloat(amountStr, 64) // 문자열을 실수형(float64)로 변환
	if err != nil {
		fmt.Println("유효하지 않은 금액입니다. 다시 시도하세요.")
		return
	}

	expenses = append(expenses, Expense{Description: desc, Amount: amount}) // 슬라이스에 새 지출 항목 추가
	fmt.Println("지출이 추가되었습니다.")
}

// 지출 목록 보기 함수
func viewExpenses() {
	if len(expenses) == 0 {
		fmt.Println("등록된 지출이 없습니다.")
		return
	}

	fmt.Println("\n지출 목록:")
	for i, expense := range expenses {
		fmt.Printf("%d. %s - %.2f\n", i+1, expense.Description, expense.Amount) // 각 항목 출력
	}
}

// CSV 파일로 내보내기 함수
func exportToCSV() {
	file, err := os.Create("expenses.csv") // 새 CSV 파일 생성
	if err != nil {
		fmt.Println("파일 생성 중 오류가 발생했습니다.")
		return
	}
	defer file.Close() // 함수 종료 시 파일 닫기

	writer := csv.NewWriter(file) // CSV 쓰기 도구 생성
	defer writer.Flush()          // 모든 데이터 쓰기 완료 후 버퍼 비우기

	for _, expense := range expenses {
		record := []string{expense.Description, fmt.Sprintf("%.2f", expense.Amount)} // CSV 형식에 맞는 레코드 생성
		if err := writer.Write(record); err != nil {
			fmt.Println("파일 쓰기 중 오류가 발생했습니다.")
			return
		}
	}
	fmt.Println("CSV 파일로 저장되었습니다: expenses.csv")
}
