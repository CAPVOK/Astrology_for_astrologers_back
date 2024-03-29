# Астрономия для астрологов

Проект создан в рамках курса "Разработка интернет-приложений" и включает в себя фронтенд, бэкенд, десктопное приложение и РПЗ. Веб-приложение выполнено в формате "Услуги/заявки", где услуги представлены планетами, а заявки — астрологическими созвездиями. Подробнее о проекте можно узнать в [РПЗ](https://github.com/CAPVOK/Astrology_for_Astrologers_documentation).

### Ссылки на репозитории проекта:
1. [Фронтенд](https://github.com/CAPVOK/Astrology_for_Astrologers_Front)
2. [Бэкенд](https://github.com/CAPVOK/Astrology_for_Astrologers_Back)
3. [Десктопное приложение](https://github.com/CAPVOK/Astrology_for_Astrologers_Desktop)
4. [GitHub Pages](https://capvok.github.io/Astrology_for_Astrologers_Front/#/)
5. [РПЗ](https://github.com/CAPVOK/Astrology_for_Astrologers_documentation)

# Бэкенд

## Ветки
- **SSR**: создание базового интерфейса, состоящего из двух страниц. Первая для просмотра списка услуг в виде карточек с наименованием и картинкой. При клике по карточке происходит переход на вторую страницу с подробной информацией об услуге. Фильтрация услуг.
- **DataBase**: разработка структуры базы данных и ее подключение к бэкенду.
- **Api**: создание веб-сервиса для получения/редактирования данных из БД, разработка всех методов для реализации итоговой бизнес-логики приложения. Соответствующая ветка фронтенда - SPA.
- **Auth**: завершение бэкенда для SPA, добавление авторизации через JWT, Swagger.

## Инструкция по запуску:
1. Клонируйте репозиторий.
2. Перейдите в директорию проекта.
3. Установите зависимости.
4. Создайте Docker: `docker-compose up -d`.
5. Выполните миграцию: `go run cmd/migration/main.go`.
6. Добавьте данные в БД.
7. Запустите приложение: `go run cmd/main.go`.

После выполнения этих шагов приложение будет доступно по адресу http://localhost:8081.

### Стек технологий:
- Go
- Gin
- Gorm
- Docker
- Minio
- Redis
- Postgres
- Swagger
