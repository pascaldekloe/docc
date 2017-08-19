.PHONY: test
test: clean code.go
	go test

.PHONY: explain
explain: clean code.dot
	xdot code.dot

code.go:
	ragel -Z code.rl

code.dot:
	ragel -V -p -o code.dot code.rl

.PHONY: clean
clean:
	rm -f code.go code.dot
