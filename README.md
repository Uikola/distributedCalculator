<h3 align="center">Yandex Distributed Calculator</h3>

  <p align="center">
    Распределённый калькулятор для вычисления выражений
    <br />
    <br />
    <a href="https://github.com/Uikola/distributedCalculator">View Demo</a>
    ·
    <a href="https://t.me/uikola">Report Bug</a>
  </p>


<!-- ABOUT TASK -->
## О Задаче
Пользователь хочет считать арифметические выражения. Он вводит строку 2 + 2 * 2 и хочет получить в ответ 6. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго. Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна. Более того: вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности. Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин". Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"? Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арфиметического выражения можно вычислять параллельно.

<!-- ABOUT THE PROJECT -->
## О Проекте
Калькулятор позволяет вычислять несколько выражений параллельно. Есть сервис оркестратор, он принимает запросы и направляет
выражения сервисам калькуляторам. Присутствуют следующие хендлеры.
- `POST /calculate` принимает задачу для вычисления и отдаёт её свободному вычислительному ресурсу(если такого нет выводит ошибку с кодом 400). В теле запроса нужно указать уникальный id задачи, чтобы в будущем узнать результат. Если задача с таким id уже существует, то отвечает кодом 200 и предупреждает о дубликате.
- `GET /tasks/{id}` возвращает задачу по её уникальному идентификатору.
- `GET /tasks` возвращает результат вычисленной задачи, если она посчиталась, иначе просит подождать.
- `GET /operations` возвращает список доступных операций с временем их выполнения.
- `PUT /operations` меняет время выполнения указанной операции.
- `POST /registry` регистрирует вычислительный ресурс.
- `GET /results/{id}` возвращает вычисленный результат по id задачи.
- `GET /c_resources` возвращает список вычислительных ресурсов с задачами, выполняющимися на них.

### Использованные Технологии

- PostgresSQL
- Apache Kafka
- Golang
- Chi router
- Docker
- Интерфейс сгенерирован с помощью swagger

<!-- GETTING STARTED -->
## Начало Работы

Чтобы запустить приложение следуйте следующим шагам.

### Установка

1. Клонируйте репозиторий
   ```sh
   git clone https://github.com/Uikola/distributedCalculator.git
   ```
2. Создайте директории envs в директориях config и добавьте туда prod.env файл(у вас должно получиться два .env файла: 1 в orchestrator/internal/config/envs, 2 в calculator/internal/config/envs).

3. Создайте директории envs в директориях config и добавьте туда prod.env файл(у вас должно получиться два .env файла: 1 в orchestrator/internal/config/envs, 2 в calculator/internal/config/envs).

4. Файл orchestrator/internal/config/envs/prod.env должен иметь следующее содержимое:
   * PORT=:8080
   * CONN_STRING=host=db port=5432 user=postgres password=password dbname=orchestratorDB sslmode=disable
   * DRIVER_NAME=postgres
   * ENV={dev or prod}
   * TIMEOUT={your_timeout}
   * IDLE_TIMEOUT={your_idle_timeout}

5. Файл calculator/internal/config/envs/prod.env должен иметь следующее содержимое:
   * CONN_STRING=host=localhost port=5432 user=postgres password=fgaSHFRdgkA4 dbname=postgres sslmode=disable
   * DRIVER_NAME=postgres
6. Соберите образ оркестратора запустив команду в его директории.
```sh
docker build -t orchestrator -f Dockerfile .
```

7. Запустите docker-compose:
 ```sh
   docker-compose up -d
   ```

8. Запустите migrator в запущенном контейнере.
 ```sh
   docker exec -it {your_container_name} go run cmd/migrator/main.go --db-url=postg
res://postgres:password@db:5432/orchestratorDB?sslmode=disable
   ```

9. Создайте топики в кафке
```sh
* docker-compose exec kafka kafka-topics.sh --create --topic expressions --partitions 1 --replication-factor 1 --z
ookeeper zookeeper:2181
* docker-compose exec kafka kafka-topics.sh --create --topic results --partitions 1 --replication-factor 1 --z
ookeeper zookeeper:2181
* docker-compose exec kafka kafka-topics.sh --create --topic heartbeat --partitions 1 --replication-factor 1 --z
ookeeper zookeeper:2181
```

10. Приложение готово к использованию! Запускайте сколько угодно калькуляторов(сперва запустити go mod tidy в директории калькулятора) go run cmd/main.go. Вы можете открыть swagger по url: localhost:8080/swagger
<!-- CONTACT -->
## Contact(Если возникли вопросы)

Yuri - [@telegram](https://t.me/uikola) - ugulaev806@yandex.ru

Project Link: [https://github.com/Uikola/distributedCalculator](https://github.com/Uikola/distributedCalculator)


