# DzShka Backend

REST API бэкенд для приложения DzShka — платформы для управления домашними заданиями по различным предметам.

## Описание

Бэкенд предоставляет API для:
- **Управления домашними заданиями** по предметам (хранение текста и файлов)
- **Загрузки файлов** (изображения, документы и т.д.)
- **Получения информации** о заданиях для каждого предмета

Приложение хранит данные в JSON файлах и поддерживает загрузку файлов на сервер.

## Поддерживаемые предметы

- Computer Graphics (computer_graphics)
- БЖД (bjd)
- Computer Practicum (com_practicum)
- IT
- English 113 (engl113)
- English 208 (engl208)
- Mathematics (math)
- ОАП (oap)
- ОСС (oss)
- ОФГ (ofg)
- ОП1С (op1c)

## Технологический стек

- **Язык**: Go 1.x
- **Стандартная библиотека**: `net/http`, `encoding/json`, `os`, `time`
- **Архитектура**: REST API с поддержкой CORS

## Установка и запуск

### Требования

- Go 1.13+
- Terminal/Command Line

### Локальный запуск

1. **Клонируйте репозиторий**:
   ```bash
   git clone https://github.com/Ryo-del/backend.git
   cd backend
   ```

2. **Скачайте зависимости** (если требуются):
   ```bash
   go mod download
   ```

3. **Запустите сервер**:
   ```bash
   go run main.go
   ```

   Сервер запустится на `http://localhost:8080`

4. **Для продакшена, соберите бинарный файл**:
   ```bash
   go build -o homework-api main.go
   ./homework-api
   ```

## API Endpoints

### Получение домашнего задания

**GET** `/api/homework/{subject}`

Получает информацию о домашнем задании для конкретного предмета.

**Пример запроса**:
```bash
curl http://localhost:8080/api/homework/math
```

**Пример ответа**:
```json
{
  "text": "Решить уравнение x² + 2x - 3 = 0",
  "files": [
    {
      "type": "image/jpeg",
      "url": "/uploads/1699000000_solution.jpg"
    }
  ],
  "updated_at": "2025-09-09 14:50"
}
```

### Создание/обновление домашнего задания

**POST** `/api/homework/{subject}`

Создает или обновляет домашнее задание для предмета.

**Пример запроса**:
```bash
curl -X POST http://localhost:8080/api/homework/math \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Решить уравнение",
    "files": [],
    "updated_at": "2025-11-27T12:00:00Z"
  }'
```

**Пример ответа**:
```json
{
  "status": "success"
}
```

### Загрузка файла

**POST** `/api/upload`

Загружает файл и возвращает URL для последующего использования.

**Параметры**:
- `file` (form-data) — файл для загрузки (макс. 10 MB)

**Пример запроса**:
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@/path/to/file.jpg"
```

**Пример ответа**:
```json
{
  "type": "image/jpeg",
  "url": "/uploads/1699000000_file.jpg"
}
```

### Получение загруженного файла

**GET** `/uploads/{filename}`

Получает загруженный файл по имени.

**Пример**:
```
http://localhost:8080/uploads/1699000000_file.jpg
```

## CORS

Сервер поддерживает CORS для всех origin'ов. Разрешены методы:
- `GET`
- `POST`
- `OPTIONS`

## Структура проекта

```
.
├── main.go           # Основной файл приложения
├── go.mod            # Go модули
├── homework.json     # Пример файла домашнего задания
├── uploads/          # Директория для загруженных файлов (создается автоматически)
└── README.md         # Этот файл
```

## Структура данных

### Homework

```go
type Homework struct {
    Text      string `json:"text"`           // Текст задания
    Files     []File `json:"files"`         // Массив файлов
    UpdatedAt string `json:"updated_at"`   // Дата обновления (формат: YYYY-MM-DD HH:mm)
}
```

### File

```go
type File struct {
    Name string `json:"name,omitempty"` // Имя файла (скрыто в ответе)
    Type string `json:"type"`           // MIME-тип файла
    URL  string `json:"url"`            // URL для доступа к файлу
}
```

## Особенности

- ✅ **CORS поддержка** — работает с фронтенд приложением dzshka
- ✅ **Загрузка файлов** — сохраняет файлы с временными метками
- ✅ **Форматирование даты** — автоматически преобразует ISO 8601 в удобный формат
- ✅ **Персистентность** — сохраняет данные в JSON файлы
- ✅ **Простота** — не требует БД, легко развертывается

## Примечания по безопасности

⚠️ **Важно**: Это учебный проект. Для продакшена рекомендуется:
- Добавить аутентификацию и авторизацию
- Валидировать входные данные
- Ограничить размер загружаемых файлов более строго
- Использовать HTTPS
- Добавить логирование
- Защитить от путем обхода директорий в загрузках
- Кэшировать результаты при необходимости

## Развертывание

### Docker

Создайте `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o homework-api main.go

EXPOSE 8080

CMD ["./homework-api"]
```

Соберите и запустите:

```bash
docker build -t homework-api .
docker run -p 8080:8080 -v $(pwd)/uploads:/app/uploads homework-api
```

## Интеграция с фронтенд (dzshka)

Фронтенд приложение (https://github.com/Ryo-del/dzshka) отправляет запросы к этому API:

```javascript
// Пример из фронтенда
fetch('http://localhost:8080/api/homework/math')
  .then(res => res.json())
  .then(data => console.log(data))
```

Убедитесь, что:
1. Бэкенд запущен на порту 8080
2. CORS включен (по умолчанию включен)
3. При изменении хоста обновите URL в фронтенд коде

## Решение проблем

### Ошибка "Address already in use"

Порт 8080 уже занят. Измените порт в `main.go`:

```go
const port = ":3000" // вместо ":8080"
```

### Файлы не загружаются

Убедитесь, что папка `uploads` имеет права на запись:

```bash
chmod 755 uploads
```

### CORS ошибки в браузере

Проверьте, что фронтенд обращается с правильным URL и методом.

## Контрибьютинг

Это личный/учебный проект. Для участия свяжитесь с автором.

## Лицензия

Личный проект, 2025

## Контакты

GitHub: [@Ryo-del](https://github.com/Ryo-del)
Фронтенд репозиторий: https://github.com/Ryo-del/dzshka
