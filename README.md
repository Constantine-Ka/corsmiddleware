# CORS Middleware
Промежуточный обработчик КОРС.
Требуется отправить сервису в POST-запросе типа formData определенные поля. 
Приложение сделает запрос, и отдаст ответ без CORS-заголовков. Смотрите примеры

## Флаги
- ``-port``: (число) ``8008`` - порт на котором будет работать веб-сервер
- ``-prometheus``:(число) ``0`` - порт, который будет слушать Prometheus. Если 0, то выключено.

## Системные поля
- ``_url`` Обязательное. URL по которому нужно выполнить запрос.
- ``_method`` http-метод, с которым нужно отправить запрос.
- ``_headers`` Строка в формате json. Oбъект с заголовками, для запроса
- ``_json`` Строка в формате json. Тело запроса, которое отправится на сервер в виде текста raw
- ``Любые поля`` Если ``_json`` пустой, метод POST, то любые поля, которые вы укажете, отправятся в запросе с типом formdata. 
Если метод GET, то эти поля добавятся к гет-запросу.

## Применение
### JSONplaceholder TODOlist
http request
```
POST / HTTP/1.1
Host: localhost:8008
Content-Length: 170
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="_url"

https://jsonplaceholder.typicode.com/todos
------WebKitFormBoundary7MA4YWxkTrZu0gW--

```
curl
```
curl --location 'localhost:8008' \
--form '_url="https://jsonplaceholder.typicode.com/todos"'
```
js
```
//Axios
let data = new FormData();
data.append('_url', 'https://jsonplaceholder.typicode.com/todos');

let config = {
    method: 'post',
    maxBodyLength: Infinity,
    url: 'localhost:8008',
    headers: {
        ...data.getHeaders()
    },
    data : data
};

axios.request(config)
    .then((response) => {
        console.log(JSON.stringify(response.data));
    })
    .catch((error) => {
        console.log(error);
    });
//Fetch
const formdata = new FormData();
formdata.append("_url", "https://jsonplaceholder.typicode.com/todos");

const requestOptions = {
    method: "POST",
    body: formdata,
    redirect: "follow"
};

fetch("localhost:8008", requestOptions)
    .then((response) => response.text())
    .then((result) => console.log(result))
    .catch((error) => console.error(error));
//jquery
var form = new FormData();
form.append("_url", "https://jsonplaceholder.typicode.com/todos");

var settings = {
    "url": "localhost:8008",
    "method": "POST",
    "timeout": 0,
    "processData": false,
    "mimeType": "multipart/form-data",
    "contentType": false,
    "data": form
};

$.ajax(settings).done(function (response) {
    console.log(response);
});

```
### Живой поиск адреса
http request
```
POST / HTTP/1.1
Host: localhost:8008
Content-Length: 526
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="_method"

POst
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="_url"

https://service.nalog.ru/static/fias-client.json
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="c"

GetAddressHint
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="query"

397907, Воронежская область, г.Лиски, ул. Василия Буракова, д. 24, кв. 50
------WebKitFormBoundary7MA4YWxkTrZu0gW--

```
curl
```curl
curl --location 'localhost:8008' \
--form '_method="POst"' \
--form '_url="https://service.nalog.ru/static/fias-client.json"' \
--form 'c="GetAddressHint"' \
--form 'query="397907, Воронежская область, г.Лиски, ул. Василия Буракова, д. 24, кв. 50"'

```
js
```javascript
//Axios
const axios = require('axios');
const FormData = require('form-data');
let data = new FormData();
data.append('_method', 'POst');
data.append('_url', 'https://service.nalog.ru/static/fias-client.json');
data.append('c', 'GetAddressHint');
data.append('query', '397907, Воронежская область, г.Лиски, ул. Василия Буракова, д. 24, кв. 50');

let config = {
    method: 'post',
    maxBodyLength: Infinity,
    url: 'localhost:8008',
    headers: {
        ...data.getHeaders()
    },
    data : data
};

axios.request(config)
    .then((response) => {
        console.log(JSON.stringify(response.data));
    })
    .catch((error) => {
        console.log(error);
    });
//Fetch
const formdata = new FormData();
formdata.append("_method", "POst");
formdata.append("_url", "https://service.nalog.ru/static/fias-client.json");
formdata.append("c", "GetAddressHint");
formdata.append("query", "397907, Воронежская область, г.Лиски, ул. Василия Буракова, д. 24, кв. 50");

const requestOptions = {
    method: "POST",
    body: formdata,
    redirect: "follow"
};

fetch("localhost:8008", requestOptions)
    .then((response) => response.text())
    .then((result) => console.log(result))
    .catch((error) => console.error(error));

//jquery
var form = new FormData();
form.append("_method", "POst");
form.append("_url", "https://service.nalog.ru/static/fias-client.json");
form.append("c", "GetAddressHint");
form.append("query", "397907, Воронежская область, г.Лиски, ул. Василия Буракова, д. 24, кв. 50");

var settings = {
    "url": "localhost:8008",
    "method": "POST",
    "timeout": 0,
    "processData": false,
    "mimeType": "multipart/form-data",
    "contentType": false,
    "data": form
};

$.ajax(settings).done(function (response) {
    console.log(response);
});
```
## Swagger
Сваггер доступен по адресу /index.html

``http://localhost:8008/index.html``

## Установка
1. Копируем бинарный файл в папку на сервере. Не забываем прописать CHMOD.
2. Переходим в папку ``/usr/lib/systemd/system``. Создаем файл заканчивающийся на .service. 
Пускай будет ``cors.service``. Прописываем внутри него директорию, путь к бинарнику, а также порты для вебсервера и прометеуса.
```
[Unit]
Description=corsMiddlleware

[Service]
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory =/var/www/vendor/bin
ExecStart=/var/www/vendor/bin/corsMiddlleware_linux_amd64 -port=8008 -prometheus=9088

[Install]
WantedBy=multi-user.target
```
3. Запускаем созданную службу ``service cors start``. Проверяем ``service cors status``
4. Заходим в nginx-конфигурацию прописываем перенаправление. Например:
```nginx
    location /cors {
		proxy_pass http://localhost:8008/;
		proxy_set_header X-Real-Ip $remote_addr;
		add_header 'Access-Control-Allow-Origin' '*';
	}
 
```
5. Проверяем конфигурацию ```nginx -t``` и перезапускаем nginx ``service nginx restart``
6. Переходим на сайт /index.html.

---

<small>Важно. Данный проект бы написан в октябре-ноябре 2024.</small>
