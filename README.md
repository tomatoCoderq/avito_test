# Avito Task - Team Management API

REST API для управления командами разработки с функциональностью управления пользователями и pull request'ами.

## Описание решения

Система предоставляет HTTP API для:
- Создания и управления командами разработки
- Добавления и деактивации пользователей в командах
- Управления pull request'ами (создание, получение, слияние)
- Переназначения reviewer'ов при деактивации пользователей
- Получения статистики по командам и пользователям

## Дополнительная функциональность

Помимо базовых требований ТЗ реализованы:

### Основные дополнения
- **addUsers endpoint** - добавление пользователей в команду
- **deactivation endpoint** - массовая деактивация пользователей с переназначением PR

### Дополнительные задачи
- **Linting** - настроен golangci-lint для проверки качества кода
- **Stress tests** - нагрузочное тестирование с k6 (100 и 300 одновременных пользователей)
- **Performance optimization** - оптимизированные SQL-запросы для массовых операций

## Архитектура

- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL с GORM ORM
- **Pattern**: Repository-Service-Controller
- **Containerization**: Docker + Docker Compose
- **Testing**: k6 для нагрузочного тестирования

## Запуск проекта

### Полный запуск через Make

```bash
# Копирование .env | Поменяйте файл под ваши нужды, если это требуется
cp .env.example .env


# Запуск всего проекта
make up

# Перезапуск сервисов
make restart

# Пересборка и запуск
make recreate

# Проверка качества кода
make lint-fix

# Остановка и очистка
make clean
```

### Нагрузочное тестирование

```bash
# Перейти в директорию тестов
cd load_tests

# Проверить готовность
make check-ready

# Запустить стресс-тесты
make stress-test

# Очистить результаты
make clean
```

### Ручной запуск

Если нужен запуск без Make:

```bash
# Запуск через Docker Compose
docker-compose up -d

# Запуск Go приложения локально
go run src/api/main.go

# Запуск тестов
cd load_tests
k6 run stress-test.js
```

## Конфигурация

Основные настройки в `docker-compose.yml`:
- API server: `localhost:8080`
- PostgreSQL: `localhost:5432`
- Database: `avito_db`
