package main

import (
	"encoding/json" // JSON 인코딩 및 디코딩을 위한 패키지
	"math/rand"     // 랜덤 숫자를 생성하기 위한 패키지
	"net/http"      // HTTP 서버 및 클라이언트 기능을 제공하는 패키지
	"strconv"       // 문자열과 숫자 간 변환을 위한 패키지

	"github.com/gorilla/mux" // 라우터를 쉽게 설정하기 위한 패키지
)

// 게시물(Post)의 구조체 정의
// 각 필드는 JSON으로 직렬화될 때의 키 이름을 지정
// ID: 게시물 ID, Title: 게시물 제목, Content: 게시물 내용
type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 메모리에 저장되는 게시물 목록 (데이터베이스 대신 사용)
var posts []Post

// 모든 게시물을 조회하는 핸들러 함수
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // 응답의 Content-Type을 JSON으로 설정
	json.NewEncoder(w).Encode(posts)                   // posts 배열을 JSON으로 변환하여 클라이언트에 응답
}

// 특정 ID의 게시물을 조회하는 핸들러 함수
func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // URL 경로에서 변수 추출
	for _, post := range posts {
		if post.ID == params["id"] { // ID가 일치하는 게시물을 찾으면
			json.NewEncoder(w).Encode(post) // JSON으로 변환하여 클라이언트에 응답
			return
		}
	}
	http.Error(w, "Post not found", http.StatusNotFound) // 게시물이 없으면 404 응답
}

// 새로운 게시물을 생성하는 핸들러 함수
func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post Post
	_ = json.NewDecoder(r.Body).Decode(&post) // 요청 본문을 디코딩하여 Post 구조체로 변환
	post.ID = strconv.Itoa(rand.Intn(100000)) // 랜덤 ID 생성 (간단한 대체)
	posts = append(posts, post)               // posts 배열에 새 게시물 추가
	json.NewEncoder(w).Encode(post)           // 생성된 게시물을 클라이언트에 응답
}

// 기존 게시물을 업데이트하는 핸들러 함수
func updatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // URL 경로에서 변수 추출
	for index, post := range posts {
		if post.ID == params["id"] { // ID가 일치하는 게시물을 찾으면
			posts = append(posts[:index], posts[index+1:]...) // 기존 게시물 삭제
			var updatedPost Post
			_ = json.NewDecoder(r.Body).Decode(&updatedPost) // 요청 본문을 디코딩하여 업데이트된 데이터로 변환
			updatedPost.ID = post.ID                         // 기존 ID 유지
			posts = append(posts, updatedPost)               // 업데이트된 게시물 추가
			json.NewEncoder(w).Encode(updatedPost)           // 업데이트된 게시물을 클라이언트에 응답
			return
		}
	}
	http.Error(w, "Post not found", http.StatusNotFound) // 게시물이 없으면 404 응답
}

// 특정 게시물을 삭제하는 핸들러 함수
func deletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // URL 경로에서 변수 추출
	for index, post := range posts {
		if post.ID == params["id"] { // ID가 일치하는 게시물을 찾으면
			posts = append(posts[:index], posts[index+1:]...) // 해당 게시물을 배열에서 삭제
			w.WriteHeader(http.StatusNoContent)               // 성공적으로 삭제했음을 나타내는 응답
			return
		}
	}
	http.Error(w, "Post not found", http.StatusNotFound) // 게시물이 없으면 404 응답
}

// 메인 함수: 서버를 시작하고 라우트를 설정
func main() {
	r := mux.NewRouter() // 새로운 라우터 생성

	// 초기 데이터 추가 (테스트용)
	posts = append(posts, Post{ID: "1", Title: "First Post", Content: "This is my first post"})
	posts = append(posts, Post{ID: "2", Title: "Second Post", Content: "This is my second post"})

	// 각 API 경로와 핸들러 함수 연결
	r.HandleFunc("/posts", getPosts).Methods("GET")           // 모든 게시물 조회
	r.HandleFunc("/posts/{id}", getPost).Methods("GET")       // 특정 게시물 조회
	r.HandleFunc("/posts", createPost).Methods("POST")        // 게시물 생성
	r.HandleFunc("/posts/{id}", updatePost).Methods("PUT")    // 게시물 업데이트
	r.HandleFunc("/posts/{id}", deletePost).Methods("DELETE") // 게시물 삭제

	http.ListenAndServe(":8000", r) // HTTP 서버 시작 (포트 8000에서 대기)
}
