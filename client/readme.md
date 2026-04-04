# curl请求方式
```shell
curl http://127.0.0.1:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-6f950310-e3fc-43fe-a26c-96fd90a192a7" \
  -d '{
        "model": "deepseek-chat",
        "messages": [
          {"role": "system", "content": "You are a helpful assistant."},
          {"role": "user", "content": "go语言是什么"}
        ],
        "stream": false
      }'
```

# go sdk
https://github.com/voocel/litellm
