package main

import (
	"encoding/json" // JSON 데이터의 인코딩 및 디코딩을 처리. 예를 들어, API 응답 데이터를 구조체로 변환하거나 구조체를 JSON으로 변환할 때 사용.
	"fmt" // 문자열 형식화 및 출력에 사용. API URL 생성 및 에러 메시지 형식화를 수행.
	"io" // 응답 본문을 읽기 위한 패키지. 예를 들어, HTTP 응답 데이터를 읽을 때 사용.
	"log" // 프로그램 실행 중 로그 메시지를 출력. 예를 들어, 치명적인 에러 발생 시 로그를 기록하고 종료.
	"net/http" // HTTP 클라이언트를 사용하여 API 요청을 보내고 응답을 처리. 예를 들어, 외부 날씨 API에 HTTP GET 요청을 보낼 때 사용.
	"os" // 운영 체제와 상호작용. 환경 변수를 가져오거나 파일 작업을 처리. API 키를 환경 변수에서 로드할 때 사용.
	"time" // 시간 및 타이머 관련 작업을 처리. Redis 캐시의 만료 시간을 설정할 때 사용.

	"github.com/go-redis/redis/v8" // Redis 클라이언트를 사용하기 위한 패키지. Redis에 데이터를 저장하거나 읽을 때 사용.
	"github.com/labstack/echo/v4" // Echo 웹 프레임워크. HTTP 요청을 처리하고 라우팅을 설정하는 데 사용.
	"golang.org/x/net/context" // Redis와 같은 비동기 작업에서 컨텍스트를 관리하기 위한 패키지.
)

type WeatherResponse struct {
	City    string  `json:"city"`    // 도시 이름을 저장.
	Temp    float64 `json:"temp"`    // 온도 정보를 저장.
	Weather string  `json:"weather"` // 날씨 상태를 저장.
}

var (
	redisClient *redis.Client // Redis 클라이언트 인스턴스.
	ctx         = context.Background() // Redis 작업에 사용할 컨텍스트.
	apiKey      string // 외부 날씨 API 키를 저장.
)

func init() {
	// 환경 변수에서 API 키를 로드. API 키가 설정되지 않은 경우 프로그램을 종료.
	apiKey = os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("환경 변수에 WEATHER_API_KEY 가 등록되지 않았습니다.") // API 키가 없으면 프로그램을 종료.
	}

	// Redis 클라이언트를 초기화. 기본적으로 localhost의 Redis 서버를 사용하도록 설정.
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis 서버의 주소.
	})
}

// getWeather 함수는 Redis 캐시에서 날씨 데이터를 확인하고, 캐시에 없으면 외부 API에서 가져와 캐시에 저장합니다.
func getWeather(city string) (*WeatherResponse, error) {
	// 1. Redis 캐시에서 데이터 검색 시도.
	cachedData, err := redisClient.Get(ctx, city).Result()
	if err == nil {
		var cachedWeather WeatherResponse
		// 캐시된 데이터를 WeatherResponse 구조체로 변환.
		if err := json.Unmarshal([]byte(cachedData), &cachedWeather); err == nil {
			return &cachedWeather, nil // 캐시된 데이터를 반환.
		}
	}

	// 2. 캐시에 데이터가 없을 경우, 외부 API에서 데이터 가져오기.
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?unitGroup=metric&include=days&key=%s&contentType=json", city, apiKey)
	fmt.Println("Request URL:", url) // 요청 URL 디버깅용 출력.
	resp, err := http.Get(url) // HTTP GET 요청.
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %v", err) // 요청 실패 시 에러 반환.
	}
	defer resp.Body.Close() // 함수 종료 시 응답 본문 닫기.

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // 응답 본문 읽기.
		return nil, fmt.Errorf("API error: status code %d, response: %s", resp.StatusCode, string(body)) // 상태 코드와 응답 본문 포함.
	}

	// 3. API 응답 데이터를 JSON으로 디코딩.
	var apiResponse map[string]interface{} // API 응답 데이터를 저장할 맵.
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err) // JSON 디코딩 실패 시 에러 반환.
	}

	// 4. 응답 데이터를 WeatherResponse 구조체에 매핑.
	weather := &WeatherResponse{
		City: city, // 요청한 도시 이름.
	}

	// 온도 데이터 가져오기 (타입 검사 포함)
	if days, ok := apiResponse["days"].([]interface{}); ok && len(days) > 0 {
		if dayData, ok := days[0].(map[string]interface{}); ok {
			if t, ok := dayData["temp"].(float64); ok {
				weather.Temp = t
			} else {
				weather.Temp = 0.0 // 기본값 설정
			}
			if w, ok := dayData["conditions"].(string); ok {
				weather.Weather = w
			} else {
				weather.Weather = "Unknown" // 기본값 설정
			}
		}
	} else {
		return nil, fmt.Errorf("invalid response structure: missing 'days' field or incorrect format")
	}

	// 5. 결과를 Redis 캐시에 저장.
	cacheData, _ := json.Marshal(weather)                    // 구조체를 JSON으로 변환.
	redisClient.Set(ctx, city, cacheData, 10*time.Minute)    // 캐시 만료 시간은 10분으로 설정.

	return weather, nil
}

// weatherHandler 함수는 HTTP 요청을 처리하며, 도시 이름을 기반으로 날씨 데이터를 반환합니다.
func weatherHandler(c echo.Context) error {
	// 1. 클라이언트 요청에서 "city" 파라미터를 가져옴.
	city := c.QueryParam("city")
	if city == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "city parameter is required"}) // 도시 파라미터가 없으면 400 에러 반환.
	}

	// 2. getWeather 함수 호출하여 날씨 데이터 가져오기.
	weather, err := getWeather(city)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()}) // 에러 발생 시 500 에러 반환.
	}

	// 3. 날씨 데이터를 JSON 형식으로 반환.
	return c.JSON(http.StatusOK, weather)
}

// main 함수는 애플리케이션의 진입점으로, 서버를 초기화하고 라우트를 설정한 후 실행합니다.
func main() {
	e := echo.New() // Echo 웹 프레임워크 인스턴스 생성.

	// 1. 라우트 설정. /weather 엔드포인트에서 weatherHandler 호출.
	e.GET("/weather", weatherHandler)

	// 2. 서버 시작. 포트 8080에서 수신 대기.
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err) // 서버 시작 실패 시 로그 출력 후 종료.
	}
}
