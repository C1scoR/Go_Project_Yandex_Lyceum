# Go_Project_Yandex_Lyceum
## **Кратко о проекте**
 - Этот проект является финальной задачей первого спринта по годовому курсу на Golang. Здесь реализован web-сервис, который под капотом имеет калькулятор, который считает простейшие арифметические выражения.

Сервис имеет один endPoint - CalcHandler, который обрабатывает все ошибки. 
Были выполнены проверки на ошибки:
- 400 Bad Status Request 
- 422 Unprocessable entity
- А также во фрагменте ниже описана ошибка 500 (да, мой калькулятор не знает что делать в такой ситуации)
```go
result, err := calculator.Calc(request.Expression)
	if err != nil {
		if errors.Is(err, calculator.ErrInvalidExpression) {
			http.Error(w, fmt.Sprintf("Invalid expression: %s", err.Error()), http.StatusBadRequest)
		} else {
			http.Error(w, "unknown error", http.StatusInternalServerError)
		}
		return
	}
```
Функция `Run` служит для того, чтобы обрабатывать запросы с терминала. 
Чтобы завершить её выполнение - нужно написать `exit`

Если что эта функция (функция Run) в main закомментирована, но всегда можно раскомментировать и потестить

## **cURL запросы**
Чтобы произвести curl (в Postman) запрос достаточно написать:
```powershell
curl --location 'localhost:8000/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
Вот для консоли:
```powershell
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2 + 2 * 2\"}" http://localhost:8000

```

Рекомендую делать это в PostMan, ниже прикрепляю скрин самой базовой ситуации:
![image](https://github.com/user-attachments/assets/aa285d6a-96e5-48ef-a772-6ca138e31a49)



### **Все тесты проходятся успешно:**
![image](https://github.com/user-attachments/assets/c155d8a6-30be-4039-91c1-dc503c2f956a)


## **Сonfig** 

Я предпочёл, что программа будет запускаться изначально на 8000 порту, но если хотите, то можете раскомментировать строчки в `ConfigFromEnv`, ну или скопировать ниже(свой .env файл я пушить не буду, там же находится супер секретный, который никому нельзя рассказывать **PORT**)
```go
func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
	config.Addr = "8000"
	}
	return config
}
```
На этом всё, вот Роналдо:
![image](https://github.com/user-attachments/assets/4109fad5-c9b7-4a52-a830-5c92dbefd932)

