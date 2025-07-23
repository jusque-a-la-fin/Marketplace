# Marketplace
## Как запустить:
```bash
git clone git@github.com:jusque-a-la-fin/Marketplace.git && cd Marketplace && docker compose up --build
```
Тесты запускаются в сервисе ['test'](https://github.com/jusque-a-la-fin/Marketplace/blob/main/compose.yaml) во время выполнения 'docker compose up --build':  
1) [Тесты на сценарий регистрации нового пользователя/авторизации зарегистрированного пользователя](https://github.com/jusque-a-la-fin/Marketplace/blob/main/internal/handlers/user/auth_test.go),  
2) [Тест на сценарий создания нового объявления](https://github.com/jusque-a-la-fin/Marketplace/blob/main/internal/handlers/user/post_test.go).  
3) [Тесты на сценарий получения ленты объявлений](https://github.com/jusque-a-la-fin/Marketplace/blob/main/internal/handlers/user/get_cards_test.go).
