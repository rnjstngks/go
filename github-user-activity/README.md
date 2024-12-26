# Github Activity 사용법

**1. 파일 복사**
```sh
git clone https://github.com/rnjstngks/go.git
```

**2. 빌드 실행**
```sh
cd go/github-user-activity
go mod init github-active.go
go build -o github-activitiy
```

빌드를 실행하고 나면, go.mod, github-activity 파일이 생성 됩니다.

**3. 동작 방법**

Github User 이름 입력
```sh
./github-activity <Github User 이름>
```