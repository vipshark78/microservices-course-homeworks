# microservices-course-homeworks
Репозиторий содержит проект из курса «Микросервисы, как в BigTech 2.0» от Олега Козырева.

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/vipshark78/6b1f6f0eaf1911e106d22bba2912187d/raw/coverage.json)

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml
