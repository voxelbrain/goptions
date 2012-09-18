`goptions` is a more complex library for command line flags

# Usage

```go
func main() {
	var options struct {
		Force bool `goptions:"--force, -f, description=Continue on error"`
		Verbosity int `goptions:"--verbose, -v, description=Level of verbosity, accumulate"`
		Wait bool `goptions:"--wait, -w, description=Wait for result, mutexgroup=behaviour"`
		Continue bool `goptions:"--continue, -c, description=Continue regardless of result, mutexgroup=behaviour"`
		Verbs {
			"create": struct {
				Name string `goptions:"--name, -n, description=Name of the new object, non-zero"`
				Obj MyObject
			},
			"delete": struct {
			}
			// "": struct {}
		}
		Help `goptions:"--help, -h, -?, description=Show this help"`
		Remainder
	}
	goptions.Parse(&options)
}
```

# TODO

* Remainder vs. Rest vs. ???
* MyObject implements Marshaller interface?
