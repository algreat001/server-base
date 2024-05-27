# Server components

## Зачем

Эта штука главный компонент для построения бэков - по сути берем этот сервер - дописываем свои контроллеры/сервисы репозитории, модели бизнес-логики
И вуаля - под капотом набор crud для
- логов
- пользователей на корпоративных сервисах (основная инфа лежит в приложении login-service)
- их права ace
- доступ к шине
- доступ к сокетам
- хелперы для построения запросов (сокеты и restful)
- хелпер и класс "процессора" - запроса с длинным циклом жизни
- менеджер "процессоров"
- хелпер для работы с файлами
- хелпер для работы с конфигами
- набор ошибок

## install for editing
```
cd existing_repo
git remote add origin https://gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server.git
git branch -M main
git push -uf origin main
```

## install for using in applications
```
export GOPRIVATE=gitlab.autocarat.de && go get gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server@latest
```
