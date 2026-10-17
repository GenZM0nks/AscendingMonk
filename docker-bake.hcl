group "default" {
  targets = ["ascendingmonk"]
}

target "ascendingmonk" {
  context    = "."
  dockerfile = "Dockerfile"
  tags       = ["ascendingmonk:local"]
}