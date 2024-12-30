# CRUD API 사용법

**1. 파일 복사**
```sh
git clone https://github.com/rnjstngks/go.git
```

**2. 빌드 실행**
```sh
cd go/blogging-platform-api
go mod init crud-api.go
go get github.com/gorilla/mux
go mod tidy
go build -o crud-api
```

빌드를 실행하고 나면, go.mod, crud-api 파일이 생성 됩니다.

**3. 동작 방법**

crud-api 실행
```sh
./crud-api
```

인터넷 브라우저 열고 난 후
-----------------------------
모든 게시물 조회
localhost:8000/posts

특정 게시물 조회
localhost:8000/posts/{id}
-----------------------------
curl 명령어를 사용하여 게시물 생성, 업데이트, 삭제

게시물 생성
curl -X PUT http://localhost:8000/posts -H "Content-Type: application/json" -d '{
  "title": "새 게시물 제목",
  "content": "새 게시물 내용"
}'

게시물 업데이트
curl -X POST http://localhost:8000/posts/{id} -H "Content-Type: application/json" -d '{
  "title": "수정된 게시물 제목",
  "content": "수정된 게시물 내용"
}'

게시물 삭제
curl -X DELETE http://localhost:8000/posts/{id}