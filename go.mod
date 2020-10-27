module github.com/Solar-2020/Group-Backend

go 1.14

require (
	github.com/Solar-2020/Authorization-Backend v0.0.0-20201027204158-15670a9b5d96
	github.com/Solar-2020/GoUtils v0.0.0-20201027194059-562c66fd0229
	github.com/buaazp/fasthttprouter v0.1.1
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.8.0
	github.com/rs/zerolog v1.20.0
	github.com/valyala/fasthttp v1.16.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

//replace github.com/Solar-2020/GoUtils => ../GoUtils

//replace github.com/Solar-2020/Authorization-Backend => ../Authorization-Backend
