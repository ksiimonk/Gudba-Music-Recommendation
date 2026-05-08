# Демо-сценарии

## Сценарий 1: Новый пользователь

**Цель:** Показать полный цикл от регистрации до персонализированных рекомендаций.

### Шаги

1. Открыть http://localhost:5173
2. Нажать "Зарегистрироваться"
3. Ввести email (например, `demo@test.local`) и пароль (`123456`)
4. Нажать "Зарегистрироваться"
5. Система автоматически логинит пользователя и редиректит на главную
6. На главной появляется CTA "Настроить рекомендации"
7. Нажать "Настроить рекомендации"
8. **Шаг 1:** Выбрать жанры "Lo-Fi Hip Hop" и "Jazz" → "Далее"
9. **Шаг 2:** Выбрать артиста "Ideal" → "Далее"
10. **Шаг 3:** Отметить 2-3 трека → "Далее"
11. **Шаг 4:** Выбрать контекст "Фокус" → "Готово"
12. Происходит редирект на главную
13. Появляется секция "Только для тебя" с персонализированными рекомендациями
14. Каждая рекомендация содержит score (зелёная полоска) и explanation ("Совпадает с твоими любимыми жанрами")

### Ожидаемый результат
- Рекомендации содержат Lo-Fi и Jazz треки
- Explanation отражает совпадение с любимыми жанрами

### Проверка через API

```powershell
.\scripts\dev.ps1 backend

# В другом терминале:
$base = "http://127.0.0.1:8080"
$creds = @{email="demo@test.local"; password="123456"} | ConvertTo-Json

# Регистрация
Invoke-RestMethod "$base/api/v1/auth/register" -Method Post -Body $creds -ContentType "application/json"

# Онбординг
$login = Invoke-RestMethod "$base/api/v1/auth/login" -Method Post -Body $creds -ContentType "application/json"
$token = $login.token
$headers = @{Authorization="Bearer $token"}
$onb = @{genre_ids=@(1,6); artist_ids=@(1); track_ids=@(); contexts=@("focus")} | ConvertTo-Json
Invoke-RestMethod "$base/api/v1/me/onboarding" -Method Post -Body $onb -ContentType "application/json" -Headers $headers

# Рекомендации
$recs = Invoke-RestMethod "$base/api/v1/recommendations/tracks?limit=5" -Headers $headers
$recs.recommendations | ForEach-Object { Write-Host "$($_.track.title) — $($_.explanation)" }
```

---

## Сценарий 2: Влияние лайков на рекомендации

**Цель:** Показать, что поведенческие события меняют рекомендации.

### Предусловие
Выполнен Сценарий 1 (пользователь зарегистрирован и прошёл онбординг).

### Шаги

1. На главной в секции "Только для тебя" найти рекомендованный трек
2. Нажать кнопку like (♥ по ховеру)
3. Нажать play на том же треке
4. Перезагрузить страницу

### Ожидаемый результат
- После like в рекомендациях может появиться explanation "Учитывает твои недавние лайки"
- Лайкнутый трек получает понижение score (−30%) в следующих рекомендациях (из-за play-события)

### Проверка через API

```powershell
# После онбординга (см. Сценарий 1)
$recs = Invoke-RestMethod "$base/api/v1/recommendations/tracks?limit=5" -Headers $headers
$trackId = $recs.recommendations[0].track.id
Write-Host "Первый рекомендованный: $($recs.recommendations[0].track.title) — $($recs.recommendations[0].score)"

# Лайк и прослушивание
Invoke-RestMethod "$base/api/v1/tracks/$trackId/like" -Method Post -Headers $headers
Invoke-RestMethod "$base/api/v1/tracks/$trackId/play" -Method Post -Headers $headers

# Повторные рекомендации
$recs2 = Invoke-RestMethod "$base/api/v1/recommendations/tracks?limit=5" -Headers $headers
Write-Host "После лайка:"
$recs2.recommendations | ForEach-Object { Write-Host "$($_.track.title) — score=$([math]::Round($_.score)) — $($_.explanation)" }

# Если трек был сыгран — его score должен понизиться
```

---

## Сценарий 3: Рекомендации плейлистов

**Цель:** Показать, что система рекомендует не только треки, но и плейлисты.

### Предусловие
Выполнен Сценарий 1 (пользователь зарегистрирован и прошёл онбординг).

### Шаги

1. Открыть главную страницу
2. После секции "Только для тебя" найти секцию "Плейлисты под твой вкус"
3. Каждый плейлист имеет explanation под карточкой

### Ожидаемый результат
- Плейлисты, содержащие треки любимых жанров/артистов, получают высокий score
- Explanation: "Совпадает с твоими любимыми жанрами" или "Содержит твоих любимых артистов"

### Проверка через API

```powershell
$plRecs = Invoke-RestMethod "$base/api/v1/recommendations/playlists?limit=6" -Headers $headers
$plRecs.recommendations | ForEach-Object { Write-Host "$($_.playlist.name) — score=$([math]::Round($_.score)) — $($_.explanation)" }

# Для пользователя, выбравшего Lo-Fi Hip Hop (genre_id=1),
# плейлист "Lo-Fi Hip Hop Mix" должен быть в топе
```

---

## Полный smoke-тест

Для автоматической проверки всех сценариев сразу выполни:

```powershell
# Терминал 1
.\scripts\dev.ps1 reset-db
.\scripts\dev.ps1 backend

# Терминал 2
.\scripts\dev.ps1 api-smoke
```
