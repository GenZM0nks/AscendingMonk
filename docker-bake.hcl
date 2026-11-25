group "default" {
  targets = ["ascendingmonk", "nginx", "playwright"]
}

target "ascendingmonk" {
  context    = "."
  dockerfile = "Dockerfile"
  tags       = ["ascendingmonk:local"]
}

target "nginx" {
  context    = "./network"
  dockerfile = "Dockerfile"
  tags       = ["ascendingmonk-nginx:local"]
}

target "playwright" {
  context = "./config/playwright"

  contexts = {
    tests = "./test/e2e"
  }

  dockerfile = "Dockerfile"
  tags = ["ascendingmonk-playwright:local"]
}