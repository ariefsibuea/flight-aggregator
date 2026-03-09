# API Documentation

## Base URL

```
http://localhost:8080
```

---

## Health Check

**`GET /health`**

Returns the service status.

**Response:**

```json
{ "status": "ok" }
```

---

## Search Flights

**`POST /api/v1/flights/search`**

Aggregates flight data from all providers, applies optional filtering and sorting, and returns a unified result set.

### Request Body

```json
{
  "origin": "CGK",
  "destination": "DPS",
  "departure_date": "2025-12-15",
  "return_date": null,
  "passengers": 1,
  "cabin_class": "economy",
  "filter": {
    "min_price": 500000,
    "max_price": 2000000,
    "max_stops": 1,
    "departure_times": ["morning", "afternoon"],
    "arrival_times": ["afternoon", "evening"],
    "airlines": ["GA", "QZ"],
    "max_duration_minutes": 200
  },
  "sort": {
    "field": "score",
    "direction": "desc"
  }
}
```

### Request Fields

| Field            | Type     | Required | Description                                                                        |
| :--------------- | :------- | :------: | :--------------------------------------------------------------------------------- |
| `origin`         | `string` |    ✅    | 3-letter IATA airport code (e.g. `"CGK"`).                                         |
| `destination`    | `string` |    ✅    | 3-letter IATA airport code (e.g. `"DPS"`).                                         |
| `departure_date` | `string` |    ✅    | Format `YYYY-MM-DD`.                                                               |
| `return_date`    | `string` |    —     | Format `YYYY-MM-DD`. Must be after `departure_date`. Triggers a round-trip search. |
| `passengers`     | `int`    |    ✅    | Number of passengers. Minimum `1`.                                                 |
| `cabin_class`    | `string` |    ✅    | One of `economy`, `business`, `first`.                                             |
| `filter`         | `object` |    —     | Optional filters applied after aggregation.                                        |
| `sort`           | `object` |    —     | Optional sort preference.                                                          |

### Filter Fields

| Field                  | Type       | Description                                                                         |
| :--------------------- | :--------- | :---------------------------------------------------------------------------------- |
| `min_price`            | `int64`    | Minimum price in IDR (inclusive).                                                   |
| `max_price`            | `int64`    | Maximum price in IDR (inclusive).                                                   |
| `max_stops`            | `int`      | Maximum number of stops. `0` = direct only.                                         |
| `departure_times`      | `[]string` | One or more named time windows: `early_morning`, `morning`, `afternoon`, `evening`. |
| `arrival_times`        | `[]string` | Same windows as `departure_times`, applied to arrival time.                         |
| `airlines`             | `[]string` | Airline codes to include (upper case). E.g. `["QZ", "GA"]`.                         |
| `max_duration_minutes` | `int`      | Maximum total trip duration in minutes (inclusive).                                 |

**Time window definitions (local airport time):**

| Window          | Range         |
| :-------------- | :------------ |
| `early_morning` | 00:00 – 06:00 |
| `morning`       | 06:00 – 12:00 |
| `afternoon`     | 12:00 – 18:00 |
| `evening`       | 18:00 – 24:00 |

Departure windows are compared against the departure airport's local time; arrival windows against the arrival airport's local time.

### Sort Fields

| Field       | Type     | Description                                                  |
| :---------- | :------- | :----------------------------------------------------------- |
| `field`     | `string` | One of `price`, `duration`, `departure`, `arrival`, `score`. |
| `direction` | `string` | `asc` or `desc`. Defaults to `asc`.                          |

---

### Response

```json
{
  "search_criteria": {
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "passengers": 1,
    "cabin_class": "economy"
  },
  "metadata": {
    "total_results": 15,
    "providers_queried": 4,
    "providers_succeeded": 4,
    "providers_failed": 0,
    "search_time_ms": 120,
    "cache_hit": false
  },
  "flights": [
    {
      "id": "GA315_GarudaIndonesia",
      "provider": "Garuda Indonesia",
      "airline": {
        "name": "Garuda Indonesia",
        "code": "GA"
      },
      "flight_number": "GA315",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T07:00:00+07:00",
        "timestamp": 1765756800
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T10:00:00+08:00",
        "timestamp": 1765764000
      },
      "duration": {
        "total_minutes": 120,
        "formatted": "2h 0m"
      },
      "stops": 1,
      "price": {
        "amount": 1850000,
        "currency": "IDR",
        "display": "Rp 1.850.000"
      },
      "amenities": ["Meal", "WiFi"],
      "baggage": {
        "carry_on": "1 pc",
        "checked": "2 pcs"
      }
    }
  ]
}
```

### Response Fields

| Field                          | Type   | Description                                       |
| :----------------------------- | :----- | :------------------------------------------------ |
| `search_criteria`              | object | Echo of the core search identity fields.          |
| `metadata.total_results`       | int    | Number of flights returned after filtering.       |
| `metadata.providers_queried`   | int    | Number of providers called. `0` on cache hit.     |
| `metadata.providers_succeeded` | int    | Number of providers that returned data.           |
| `metadata.providers_failed`    | int    | Number of providers that failed or timed out.     |
| `metadata.search_time_ms`      | int64  | Wall-clock time for the search in milliseconds.   |
| `metadata.cache_hit`           | bool   | `true` if the result was served from Redis cache. |
| `flights`                      | array  | Normalized flight results (one-way).              |
| `outbound_flights`             | array  | Round-trip only: outbound leg.                    |
| `return_flights`               | array  | Round-trip only: return leg.                      |

**Flight object:**

| Field                           | Description                                                                     |
| :------------------------------ | :------------------------------------------------------------------------------ |
| `id`                            | Unique identifier: `{flight_number}_{provider}`.                                |
| `provider`                      | Source provider name.                                                           |
| `airline.name` / `airline.code` | Airline display name and IATA code.                                             |
| `flight_number`                 | Airline flight number (e.g. `"GA315"`).                                         |
| `departure` / `arrival`         | Airport code, city, ISO 8601 datetime with timezone offset, and Unix timestamp. |
| `duration.total_minutes`        | Total trip duration in minutes (includes layovers for connecting flights).      |
| `duration.formatted`            | Human-readable string, e.g. `"2h 0m"`.                                          |
| `stops`                         | Number of stops. `0` = direct flight.                                           |
| `price.amount`                  | Raw price in IDR.                                                               |
| `price.currency`                | Always `"IDR"`.                                                                 |
| `price.display`                 | Formatted price string, e.g. `"Rp 1.850.000"`.                                  |
| `amenities`                     | List of available amenities, e.g. `["Meal", "WiFi", "Entertainment"]`.          |
| `baggage.carry_on`              | Carry-on baggage allowance.                                                     |
| `baggage.checked`               | Checked baggage allowance.                                                      |

---

### Error Responses

All errors return a structured JSON body with an HTTP `400` status:

```json
{
  "error": "validation failed",
  "details": "origin is required"
}
```

| HTTP Status                 | When                                                                                                                           |
| :-------------------------- | :----------------------------------------------------------------------------------------------------------------------------- |
| `400 Bad Request`           | Missing required fields, invalid date format, invalid `sort.field` or `sort.direction`, `return_date` before `departure_date`. |
| `500 Internal Server Error` | Unexpected server-side failure.                                                                                                |

> Provider failures do **not** produce a non-200 response. If one or more providers fail, the response still returns `200 OK` with results from the remaining providers and the failure count noted in `metadata.providers_failed`.
