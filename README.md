# Temperature System

Serviço HTTP em Go que recebe um CEP (8 dígitos), identifica a cidade via ViaCEP e retorna a temperatura atual consultando a WeatherAPI, com conversões para Celsius, Fahrenheit e Kelvin.

---

## URL de produção

* [https://weather-service-531929244051.southamerica-east1.run.app/](https://weather-service-531929244051.southamerica-east1.run.app/)

## API

### `GET /weather?cep=<CEP_8_DIGITOS>`

**Query param**

* `cep` — obrigatório; somente dígitos (aceita formatos com hífen na entrada, que são normalizados).

**Exemplos de chamada**

*Localhost*

```bash
curl "http://localhost:8080/weather?cep=01001000"
```

*Produção*

````bash
curl "https://weather-service-531929244051.southamerica-east1.run.app/weather?cep=01001000"
````

---

## Pré‑requisitos

* Go 1.20+
* Conta na **WeatherAPI** (chave em `WEATHER_API_KEY`)

---

## Configuração

Variáveis de ambiente importantes:

* `WEATHER_API_KEY` — **obrigatória** para consultar a WeatherAPI.

---

## Executar localmente

```bash
export WEATHER_API_KEY="SUA_CHAVE_WEATHERAPI"
go run .

# Teste
curl "http://localhost:8080/weather?cep=01001000"
```

* Para a WEATHER_API_KEY, você também pode criar um arquivo .env na raiz do projeto.

---

## Docker

```bash
# build & up
WEATHER_API_KEY="SUA_CHAVE_WEATHERAPI" docker compose up --build

# teste
curl "http://localhost:8080/weather?cep=01001000"
```

---

## Testes automatizados

Execute todos os testes:

```bash
go test ./... -v
```
