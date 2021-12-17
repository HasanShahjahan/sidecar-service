### About the task ###

fxproxy is sidecar service which is responsible to govern incoming traffic to the downstream application services based on allowed path list.

## Things you need

* [Go](https://golang.org/dl/): `brew install go`
* Setup `go` folder in your home directory: `/Users/<user_id>/go/src/github.com/sidecar-service`. So, `go env` should show `GOPATH="/Users/<user_id>/go"`
* Clone `sidecar-service` repository in `/Users/<user_id>/go/src/github.com/sidecar-service`: `git clone <repo_url>`
* After cloning, the `go.mod` file should be found in this directory `/Users/<user_id>/go/src/github.com/sidecar-service/go.mod`
* Add `/Users/<user_id>/go/bin` in environment variable `PATH` if it is not already there.

## Running the App

### Run Scripts
We use [Go Modules](https://blog.golang.org/using-go-modules) for dependency management.

### MAKEFILE IS YOUR FRIEND, ANY COMMAND YOU FIND BELOW IS ALREADY THERE , WITH HELP
1. `go mod vendor` should create a `vendor` directory in project root directory. (`/Users/<user_id>/go/src/github.com/sidecar-service/vendor`)
   If for any reason, this directory is not created then try to clear the cache and run the command again. To clear the cache run below command:

```shell
$ go clean -modcache  #clean module cache
$ go mod vendor       #setup vendor dir again
```