## Тестовое задание на позицию стажера backend в юнит Geo
### Quick start:
    docker-compose up   
### API methods:  
#### /login GET
Для авторизации необходимо передать Query параметр, где key=password, а
value=supersecretpassword(пароль можно изменить в configs/config.json).
После выполнение этого запроса, будет установлена cookie, открывающая вам
доступ к остальным запросам
   
    http://localhost:8080/login?password=supersecretpassword
#### /save GET
Сохранит в файл data.txt все имеющиеся данные. Чтобы просмотреть файл, можно
написать следующую команду в docker-cli

    # cat data.txt
    verygoodkey55 - verygoodvalue110
    verygoodkey555 - verygoodvalue110
    verygoodkey5 - verygoodvalue110
#### /keys GET
Позволяет получить все существующие ключи

    {
        "Keys": [
            "verygoodkey2",
            "verygoodkey5",
            "verygoodkey55"
        ]
    }
#### /rows/{key} GET
Возвращает значение по ключу

    http://localhost:8080/rows/verygoodkey2
    
#### /rows POST
Создает новую строку с ключом, значением и временем жизни.Для бессмертия,
задайте ttl=0

    http://localhost:8080/rows?key=verygoodkey55&value=verygoodvalue110&ttl=10
#### /rows/{key} DELETE
Удаляет по ключу

    http://localhost:8080/rows/verygoodkey2
