# Weather-api 사용법

**1. 파일 복사**
```sh
git clone https://github.com/rnjstngks/go.git
```

**2. 빌드 실행**
```sh
cd go/weather-api
go mod init weather-api.go
go get github.com/go-redis/redis/v8
go get github.com/labstack/echo/v4
go get golang.org/x/net/context
go mod tidy
go build -o weather-api
```

빌드를 실행하고 나면, go.mod, weather-api 파일이 생성 됩니다.

**3. 동작 방법**

환경 변수 입력
```sh
export WEATHER_API_KEY=<API_KEY 입력>
```

redis 설치
```sh
apt install -y redis
```

바이너리 파일 실행
```sh
./weather-api
```

인터넷 브라우저 열고 난 후

localhost:8080/weather?city=<원하는 도시 명 입력>