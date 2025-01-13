package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4                                                     // Устанавливаем ожидаемое количество кафе
	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil) // Создаем запрос с параметрами count=5 и city=moscow

	responseRecorder := httptest.NewRecorder() // Создаем запись ответа
	handler := http.HandlerFunc(mainHandle)    // Создаем обработчик
	handler.ServeHTTP(responseRecorder, req)   // Отправляем запрос обработчику

	require.Equal(t, http.StatusOK, responseRecorder.Code)                        // Проверяем, что код ответа равен 200
	assert.NotEmpty(t, responseRecorder.Body.String())                            // Проверяем, что тело ответа не пустое
	assert.Len(t, strings.Split(responseRecorder.Body.String(), ","), totalCount) // Проверяем, что количество кафе в ответе равно ожидаемому
}

func TestMainHandlerWhenCityIsWrong(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=1&city=wrongcity", nil) // Создаем запрос с параметрами count=1 и city=wrongcity

	responseRecorder := httptest.NewRecorder() // Создаем запись ответа
	handler := http.HandlerFunc(mainHandle)    // Создаем обработчик
	handler.ServeHTTP(responseRecorder, req)   // Отправляем запрос обработчику

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)      // Проверяем, что код ответа равен 400
	assert.Equal(t, "wrong city value", responseRecorder.Body.String()) // Проверяем, что тело ответа содержит сообщение "wrong city value"
}

func TestMainHandlerWhenRequestIsCorrect(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil) // Создаем запрос с параметрами count=2 и city=moscow

	responseRecorder := httptest.NewRecorder() // Создаем запись ответа
	handler := http.HandlerFunc(mainHandle)    // Создаем обработчик
	handler.ServeHTTP(responseRecorder, req)   // Отправляем запрос обработчику

	require.Equal(t, http.StatusOK, responseRecorder.Code)               // Проверяем, что код ответа равен 200
	assert.NotEmpty(t, responseRecorder.Body.String())                   // Проверяем, что тело ответа не пустое
	assert.Len(t, strings.Split(responseRecorder.Body.String(), ","), 2) // Проверяем, что количество кафе в ответе равно 2
}
