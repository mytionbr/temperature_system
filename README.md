# Weather by CEP (Go + Cloud Run)

Serviço HTTP em Go que recebe um CEP (8 dígitos), identifica a cidade via ViaCEP e retorna a temperatura atual consultando a WeatherAPI, com conversões para Celsius, Fahrenheit e Kelvin.

---

## URL de produção

* **Base**: [https://weather-service-531929244051.southamerica-east1.run.app/](https://weather-service-531929244051.southamerica-east1.run.app/)
* **Exemplo**: [https://weather-service-531929244051.southamerica-east1.run.app/weather?cep=01001000](https://weather-service-531929244051.southamerica-east1.run.app/weather?cep=01001000)

## API

### `GET /weather?cep=<CEP_8_DIGITOS>`

**Query param**

* `cep` — obrigatório; somente dígitos (aceita formatos com hífen na entrada, que são normalizados).

**Responses**

* **200 OK**

  ```json
  { "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }
  ```
* **422 Unprocessable Entity** (CEP com formato inválido)

  ```
  invalid zipcode
  ```
* **404 Not Found** (CEP não encontrado no ViaCEP)

  ```
  can not find zipcode
  ```
* **500 Internal Server Error**

**Exemplos de chamada**

*Localhost*

```bash
curl "http://localhost:8080/weather?cep=01001000"
```

*Produção*

````bash
curl "https://weather-service-531929244051.southamerica-east1.run.app/weather?cep=01001000"
```bash
curl "http://localhost:8080/weather?cep=01001000"
````

---

## Pré‑requisitos

* Go 1.20+ (recomendado 1.22+)
* Conta na **WeatherAPI** (chave em `WEATHER_API_KEY`)

---

## Configuração

Variáveis de ambiente importantes:

* `WEATHER_API_KEY` — **obrigatória** para consultar a WeatherAPI.
* `PORT` — porta do servidor HTTP (padrão `8080`).

---

## Executar localmente

```bash
export WEATHER_API_KEY="SUA_CHAVE_WEATHERAPI"
go run .

# Teste
curl "http://localhost:8080/weather?cep=01001000"
```

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
