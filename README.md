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
    
#### /rows POST
Создает новую строку с ключом, значением и временем жизни.Для бессмертия,
задайте ttl=0

    http://localhost:8080/rows?key=verygoodkey55&value=verygoodvalue110&ttl=10

#### /rows/{key} GET
Возвращает значение по ключу

    http://localhost:8080/rows/verygoodkey2
    
#### /hash POST
Создает новую строку, в качестве значения который хэш-таблица

    http://localhost:8080/rows/hash?hash=myhash3&field=myfield2&value=myvalue222

#### /hash/{hashName}/{field} GET
Возвращает значение поля в хэш таблице

    http://localhost:8080/rows/hash/myhash3/myfield2

#### /keys GET
Позволяет получить все существующие ключи

    http://localhost:8080/keys?pattern=*
    
Поддерживаемые паттерны:
+ ul  h?llo matches hello, hallo and hxllo
+ h*llo matches hllo and heeeello
+ h[ae]llo matches hello and hallo, but not hillo
+ h[^e]llo matches hallo, hbllo, ... but not hello
+ h[a-b]llo matches hallo and hbllo

#### /rows/{key} DELETE
Удаляет по ключу

    http://localhost:8080/rows/verygoodkey2

#### /save GET
Сохранит в файл data.txt все имеющиеся данные. Чтобы просмотреть файл, можно
написать следующую команду в docker-cli

    # cat data.txt
    verygoodkey55 - verygoodvalue110
    verygoodkey555 - verygoodvalue110
    verygoodkey5 - verygoodvalue110
    
