.PHONY: test
test: clean c.go
	go test

.PHONY: explain
explain: clean c.dot
	xdot c.dot

c.go:
	ragel -Z c.rl

c.dot:
	ragel -V -p -o c.dot c.rl

.PHONY: clean
clean:
	rm -f c.go c.dot
