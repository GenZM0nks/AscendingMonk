# What to do before pushing
- Generate OpenAPI spec if you wrote any new endpoints.
- Lint the workspace
- Run all tests

### Generating OpenAPI spec
Make sure you're in the root directory and that you've installed [swaggo](https://github.com/swaggo/swag). Then run
```bash
swag init -g cmd/main/main.go -o docs/
```

### Linting 
Make sure you've installed [revive](https://revive.run/). Then do `CTRL + SHIFT + p` and select `Go: Lint workspace`. The linter is configured automatically if you use VS Code. It is also configured to lint the workspace when you save a file.