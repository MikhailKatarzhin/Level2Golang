package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

func main() {
	server, err := initNewServer()
	if err != nil {
		log.Fatalf("Сбой запуска сервера: %s", err)
		os.Exit(1)
	}

	log.Printf("Сервер запущен на порте %s", server.Port)

	log.Fatal(http.ListenAndServe("localhost:"+server.Port, server.LoggingMiddleware(http.DefaultServeMux)))
}

/*
	  = == ==         == == =
	= ==== БИЗНЕС-ЛОГИКА ==== =
	  = == ==         == == =
*/

// === Структуры ===

type Event struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// === EventStore (хранилище событий в памяти) ===
type EventStore struct {
	sync.Mutex
	Events map[int]Event
	NextID int
}

// InitNewEventStore возвращает указатель на новую структуру EventStore с сначальным значением NextID = 1
func InitNewEventStore() *EventStore {
	return &EventStore{
		Events: make(map[int]Event),
		NextID: 1,
	}
}

// AddEvent добавляет событие в хранилище
func (es *EventStore) AddEvent(event Event) (int, error) {

	if event.UserID < 1 {
		return -1, fmt.Errorf("должен быть указан UserID превосходящий 0")
	}

	if event.Title == "" {
		return -1, fmt.Errorf("обязателен к заполнению Title события (!= \"\") ")
	}

	event.CreatedAt = time.Now()
	event.UpdatedAt = event.CreatedAt

	es.Lock()
	defer es.Unlock()

	event.ID = es.NextID
	es.NextID++

	es.Events[event.ID] = event

	return event.ID, nil
}

// UpdateEvent обновляет событие в хранилище
func (es *EventStore) UpdateEvent(event Event) error {
	es.Lock()
	defer es.Unlock()
	if e, exists := es.Events[event.ID]; !exists {
		return fmt.Errorf("событие с ID %d не найдено", event.ID)
	} else {
		event.CreatedAt = e.CreatedAt
		event.UpdatedAt = time.Now()

		if event.UserID < 1 {
			event.UserID = e.UserID
		}

		if event.Title == "" {
			event.Title = e.Title
		}

		if tmpTime := time.Date(1, 1, 1, 0, 0, 0, 0, e.Date.Location()); tmpTime == event.Date {
			event.Date = e.Date
		}

		es.Events[event.ID] = event
		return nil
	}
}

// DeleteEvent удаляет событие из хранилища
func (s *EventStore) DeleteEvent(eventID int) error {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.Events[eventID]; !exists {
		return fmt.Errorf("событие с ID %d не найдено", eventID)
	}

	delete(s.Events, eventID)

	return nil
}

// GetEventsByDate возвращает все события за определенную дату
func (s *EventStore) GetEventsByDate(date time.Time) []Event {
	s.Lock()
	defer s.Unlock()

	var result []Event
	for _, event := range s.Events {
		if event.Date.Year() == date.Year() && event.Date.YearDay() == date.YearDay() {
			result = append(result, event)
		}
	}

	return result
}

// GetEventsForRange возвращает все события за указанный диапазон дат
func (s *EventStore) GetEventsForRange(start, end time.Time) []Event {
	s.Lock()
	defer s.Unlock()

	var result []Event
	for _, event := range s.Events {
		if !event.Date.Before(start) && !event.Date.After(end) {
			result = append(result, event)
		}
	}

	return result
}

/*
	  = == ==       == == =
	= ==== HTTP СЕРВЕР ==== =
	  = == ==       == == =
*/

// === Структуры ===

type Server struct {
	Port     string      `json:"port"`
	Calendar *EventStore `json:"calendar"`
}

type Config struct {
	Port string `json:"port"`
}

type RequestObjects struct {
	User_ID int       `json:"user_id"`
	Date    time.Time `json:"date"`
}

// === Конфигурации ===

func LoadConfig() (Config, error) {
	config, err := ReadConfigFromFile("config.json")

	if err != nil {
		return Config{}, err
	}
	return config, nil
}

// LoadConfig читает конфиг из конфигурационного файла
func ReadConfigFromFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("ошибка чтения конфигурационного файла: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("ошибка парсинга содежимого конфигурационного файла: %v", err)
	}

	if config.Port == "" {
		return Config{}, fmt.Errorf("в конфигурационном файле не указан порт: %v", err)
	}

	return config, nil
}

// SetupRoutes задаёт систему маршрутищации
func (s *Server) SetupRoutes() {
	http.HandleFunc("/create_event", s.CreateEventHandler)
	http.HandleFunc("/update_event", s.UpdateEventHandler)
	http.HandleFunc("/delete_event", s.DeleteEventHandler)

	http.HandleFunc("/events_for_day", s.EventsForDayHandler)
	http.HandleFunc("/events_for_week", s.EventsForWeekHandler)
	http.HandleFunc("/events_for_month", s.EventsForMonthHandler)
}

func initNewServer() (*Server, error) {
	config, err := LoadConfig() // Порт нужно брать из конфигурационного файла
	if err != nil {
		return nil, fmt.Errorf("невозможно загрузить конфиг сервера: %s", err)
	}

	server := &Server{
		Port:     config.Port,
		Calendar: InitNewEventStore(),
	}

	server.SetupRoutes()

	return server, nil
}

// === Вспомогательные функции для сериализации объектов доменной области в JSON ===

// RespondWithJSON сериализует payload в JSON и отправляет его клиенту
func (s *Server) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Fatalf("Ошибка при попытке записи ответа в соединение")
	}
}

// === Вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event ===

// DecodeJSONBody парсит тело запроса JSON в структуру.
func (s *Server) DecodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

// ValidateDate проверяет корректность формата даты
func (s *Server) ValidateDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("формат даты должен соответствовать шаблону гггг-мм-дд")
	}
	return date, nil
}

// ValidateUserID проверяет корректность user_id
func (s *Server) ValidateUserID(userIDStr string) (int, error) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		return 0, fmt.Errorf("некорректный user_id")
	}
	return userID, nil
}

// ParseRequestToRequestObjects Проверяет корректность id и даты. В случае успеха возвращает структуру с извлечёнными объектами.
func (s *Server) ParseRequestToRequestObjects(r *http.Request) (RequestObjects, error) {
	userIDStr := r.URL.Query().Get("user_id")

	userID, err := s.ValidateUserID(userIDStr)
	if err != nil {
		return RequestObjects{}, fmt.Errorf("валидация пользовательского идентификатора не пройдена: %v", err)
	}

	dateStr := r.URL.Query().Get("date")

	date, err := s.ValidateDate(dateStr)
	if err != nil {
		return RequestObjects{}, fmt.Errorf("валидация даты не пройдена: %v", err)
	}

	return RequestObjects{
		User_ID: userID,
		Date:    date,
	}, nil
}

// === HTTP обработчики для каждого из методов API ===

// CreateEventHandler обрабатывает запрос на создание события
func (s *Server) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный http метод"})
		return
	}

	var event Event
	err := s.DecodeJSONBody(r, &event)
	if err != nil {
		s.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Ошибка в процессе парсинга нового события: %v", err)})
		return
	}

	eventID, err := s.Calendar.AddEvent(event)
	if err != nil {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка в процессе добавления нового события: %v", err)})
		return
	}

	s.RespondWithJSON(w, http.StatusOK, map[string]string{"result": fmt.Sprintf("Событие [ID:%d] успешно создано", eventID)})
}

// UpdateEventHandler обрабатывает запрос на обновление события
func (s *Server) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный http метод"})
		return
	}

	var event Event
	err := s.DecodeJSONBody(r, &event)
	if err != nil {
		s.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Ошибка в ходе парсинга входных параметров: %v", err)})
		return
	}

	event.UpdatedAt = time.Now()

	err = s.Calendar.UpdateEvent(event)
	if err != nil {
		s.RespondWithJSON(w, http.StatusServiceUnavailable, map[string]string{"error": fmt.Sprintf("Ошибка в процессе обновления события [ID:%d]: %v", event.ID, err)})
		return
	}

	s.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "обновление события успешно"})
}

// DeleteEventHandler обрабатывает запрос на удаление события
func (s *Server) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный http метод"})
		return
	}

	var event Event
	err := s.DecodeJSONBody(r, &event)
	if err != nil {
		s.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Ошибка в ходе парсинга входных параметров: %v", err)})
		return
	}

	if event.ID < 1 {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "некорректный ID события"})
		return
	}

	err = s.Calendar.DeleteEvent(event.ID)

	if err != nil {
		s.RespondWithJSON(w, http.StatusServiceUnavailable, map[string]string{"error": fmt.Sprintf("Ошибка в процессе удаления события [ID:%d]: %v", event.ID, err)})
		return
	}

	s.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "удаление события успешно"})
}

// EventsForDayHandler возвращает все события за конкретный день
func (s *Server) EventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный http метод"})
		return
	}

	requestObjects, err := s.ParseRequestToRequestObjects(r)
	if err != nil {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка в процессе парсинга объектов домена: %v", err)})
		return
	}

	events := s.Calendar.GetEventsByDate(requestObjects.Date)

	s.RespondWithJSON(w, http.StatusOK, events)
}

// EventsForWeekHandler возвращает события за конкретную неделю
func (s *Server) EventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный http метод"})
		return
	}

	requestObjects, err := s.ParseRequestToRequestObjects(r)
	if err != nil {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка в процессе парсинга объектов домена: %v", err)})
		return
	}

	weekday := requestObjects.Date.Weekday()

	// Определяем смещение от текущей даты до понедельника соответствующей недели
	offsetToMonday := int(time.Monday - weekday)
	if offsetToMonday > 0 {
		offsetToMonday -= 7
	}

	start := requestObjects.Date.AddDate(0, 0, offsetToMonday)

	end := start.AddDate(0, 0, 6)

	events := s.Calendar.GetEventsForRange(start, end)

	s.RespondWithJSON(w, http.StatusOK, events)
}

// EventsForMonthHandler возвращает события за конкретный месяц
func (s *Server) EventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный http метод"})
		return
	}

	requestObjects, err := s.ParseRequestToRequestObjects(r)
	if err != nil {
		s.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка в процессе парсинга объектов домена: %v", err)})
		return
	}

	start := time.Date(requestObjects.Date.Year(), requestObjects.Date.Month(), 1, 0, 0, 0, 0, requestObjects.Date.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	events := s.Calendar.GetEventsForRange(start, end)

	s.RespondWithJSON(w, http.StatusOK, events)
}

// === Middleware для логирования запросов ===

// LoggingMiddleware логирует каждый запрос
func (s *Server) LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("[%s] %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}
