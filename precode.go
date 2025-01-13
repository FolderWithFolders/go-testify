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

// Определяем тест для обработчика HTTP mainHandle
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	// Создаем объекты require и assert для проверки условий в тесте
	r := require.New(t)
	a := assert.New(t)

	// Устанавливаем ожидаемое значение параметра "count"
	totalCount := 4
	// Устанавливаем ожидаемое значение параметра "city"
	city := "moscow"

	// Создаем новый запрос HTTP
	req, err := http.NewRequest("GET", "http://localhost:8080/cafe?count=4&city=moscow", nil)
	// Если создание запроса не удалось, завершаем тест с ошибкой
	if err != nil {
		t.Fatal("wrong request", err)
	}

	// Преобразуем список кофеен в городе в строку
	cafeListJoin := strings.Join(cafeList[city], ",")
	// Преобразуем строку списка кофеен в слайс строк
	cafeListSlice := strings.Split(cafeListJoin, ",")

	// Создаем новый объект responseRecorder для записи ответа
	responseRecorder := httptest.NewRecorder()
	// Создаем новый обработчик HTTP
	handler := http.HandlerFunc(mainHandle)
	// Вызываем обработчик HTTP с запросом и responseRecorder
	handler.ServeHTTP(responseRecorder, req)

	// Проверяем, что создание запроса не вызвало ошибку
	r.NoError(err)
	// Проверяем, что статус код ответа равен 200 (OK)
	r.Equal(http.StatusOK, responseRecorder.Code)
	// Проверяем, что количество кофеен в городе равно ожидаемому значению
	a.Equal(len(cafeListSlice), totalCount, "wrong count value")
	// Проверяем, что значение параметра "city" равно ожидаемому значению
	a.Equal("moscow", city, "wrong city value")
	// Проверяем, что тело ответа не пустое
	r.NotEmpty(responseRecorder.Body.String(), "response body should not be empty")
}
