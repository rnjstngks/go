# Expense Tracker 사용법

**1. 파일 복사**
```sh
git clone https://github.com/rnjstngks/go.git
```

**2. 빌드 실행**
```sh
cd go/expense-tracker
go mod init expense-tracker.go
go build -o expense-tracker
```

빌드를 실행하고 나면, go.mod, expense-tracker 파일이 생성 됩니다.

**3. 동작 방법**

```sh
./expense-tracker
```

위 명령어를 실행하면 아래와 같이 나오고 원하는 번호를 입력해서 진행하면 됩니다.

----------------------------------------------------------------------------------------------------------------------------

지출 관리 프로그램
1. 지출 추가
2. 지출 목록 보기
3. CSV 파일로 내보내기
4. 종료
선택: