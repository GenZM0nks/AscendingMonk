group "default" {
  targets = ["ascendingmonk", "nginx"]
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
