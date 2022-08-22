# Генератор гистов для студентов

Скрипт генерирует гист по шаблону. Нужно передать только задания: 
22 Августа -- 28 Августа
1. Сделать первое задание
2. Сделать Второе задание

## How to set up

1. Install Go https://go.dev/doc/install
2. Clone Repo 
3. get github token https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
4. cd to repo folder
5. run `go build -ldflags "-X main.token=<GITHUB_TOKEN>"`
6. ./ggist "link/to/profile" "plan"
`

## Json body 
```json
// https://api.github.com/gists
"description": "План для ученика"
"public": false,
"files": {
    "plan.md": {
        "content": "task list..."
    }
}
```
