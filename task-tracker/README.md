# Task Tracker 사용법

**1. 파일 복사**
```sh
git clone https://github.com/rnjstngks/go.git
```

**2. 빌드 실행**
```sh
cd go/task-tracker
go mod init task-tracker.go
go build -o task-tracker
```

빌드를 실행하고 나면, go.mod, task-tracker 파일이 생성 됩니다.

**3. 동작 방법**

작업 추가
```sh
./task-tracker add "add work"
```

작업 수정
```sh
./task-tracker update <작업 ID> "edit work"
```

작업 삭제
```sh
./task-tracker delete <작업 ID>
```

작업 상태 변경 (진행 중 일 떄,)
```sh
./task-tracker mark-in-progress <작업 ID>
```

작업 상태 변경 (완료 상태)
```sh
./task-tracker mark-done <작업 ID>
```

작업 목록 조회 (전체)
```sh
./task-tracker list
```