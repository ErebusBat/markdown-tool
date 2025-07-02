# Phone Number Examples

This document shows examples of phone number inputs and their expected markdown transformations.

# Input Formats

## Seven Digits

Matches only on 7 digits, in specifc order or with specific seperators.

| Input      | Notes    | Output                         |
|------------|----------|--------------------------------|
| `1234567`  | Match    | `📞 [123-4567](tel:1234567)` |
| `123-4567` | Match    | `📞 [123-4567](tel:1234567)` |
| `123.4567` | Match    | `📞 [123-4567](tel:1234567)` |
| `123 4567` | No Match | `123 4567`                |
| `123,4567` | No Match | `123,4567`                |
| `01234567` | No Match | `01234567`                |

### Ten Digits

Matches only on 10 digits, in specifc order or with specific seperators.

| Input            | Notes    | Output                                |
|------------------|----------|---------------------------------------|
| `8901234567`     | Match    | `📞 [890-123-4567](tel:8901234567)` |
| `890-123-4567`   | Match    | `📞 [890-123-4567](tel:8901234567)` |
| `890.123.4567`   | Match    | `📞 [890-123-4567](tel:8901234567)` |
| `(890) 123-4567` | Match    | `📞 [890-123-4567](tel:8901234567)` |
| `(890)123-4567`  | Match    | `📞 [890-123-4567](tel:8901234567)` |
| `(890)1234567`   | Match    | `📞 [890-123-4567](tel:8901234567)` |
| `89012345670`    | No Match | `89012345670`                    |
| `890 123 4567`   | No Match | `890 123 4567`                   |
| `(890) 123 4567` | No Match | `(890) 123 4567`                 |
| `(890) 1234 567` | No Match | `(890) 1234 567`                 |
| `(890)123-456`   | No Match | `(890)123-456`                   |
| `(890)12345679`  | No Match | `(890)12345679`                  |

### Ten Digits Plus Country Code

Matches only on 10 digits, plus a 1 digit country code, in specifc order or with specific seperators.

- If the country code is `1` then it is not required to be prefixed with `+`
- Otherwise all other country codes must be prefixed with `+`

| Input                | Notes | Output                                       |
|----------------------|-------|----------------------------------------------|
| `18901234567`        | Match | `📞 [1-890-123-4567](tel:+18901234567)`    |
| `1-890-123-4567`     | Match | `📞 [1-890-123-4567](tel:+18901234567)`    |
| `1.890.123.4567`     | Match | `📞 [1-890-123-4567](tel:+18901234567)`    |
| `1 (890) 123-4567`   | Match | `📞 [1-890-123-4567](tel:+18901234567)`    |
| `1(890)123-4567`     | Match | `📞 [1-890-123-4567](tel:+18901234567)`    |
| `1(890)1234567`      | Match | `📞 [1-890-123-4567](tel:+18901234567)`    |
| `+78901234567`       | Match | `📞 [+7-890-123-4567](tel:+78901234567)`   |
| `+7-890-123-4567`    | Match | `📞 [+7-890-123-4567](tel:+78901234567)`   |
| `+7.890.123.4567`    | Match | `📞 [+7-890-123-4567](tel:+78901234567)`   |
| `+7 (890) 123-4567`  | Match | `📞 [+7-890-123-4567](tel:+78901234567)`   |
| `+7(890)123-4567`    | Match | `📞 [+7-890-123-4567](tel:+78901234567)`   |
| `+7(890)1234567`     | Match | `📞 [+7-890-123-4567](tel:+78901234567)`   |