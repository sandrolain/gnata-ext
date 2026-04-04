// Package main demonstrates extcrypto functions with gnata.
package main

import (
	"context"
	"fmt"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extcrypto"
)

func main() {
	env := gnata.NewCustomEnv(extcrypto.All())
	ctx := context.Background()

	eval := func(expr string, data any) any {
		e, err := gnata.Compile(expr)
		if err != nil {
			panic(fmt.Sprintf("compile %q: %v", expr, err))
		}
		result, err := e.EvalWithCustomFuncs(ctx, data, env)
		if err != nil {
			panic(fmt.Sprintf("eval %q: %v", expr, err))
		}
		return result
	}

	// $uuid() – generate a random UUID v4
	fmt.Println("uuid:", eval(`$uuid()`, nil))

	// $hash(algorithm, value) – hex-encoded digest
	// Supported algorithms: md5, sha1, sha256, sha384, sha512
	fmt.Println("md5:   ", eval(`$hash("md5", "hello")`, nil))
	fmt.Println("sha1:  ", eval(`$hash("sha1", "hello")`, nil))
	fmt.Println("sha256:", eval(`$hash("sha256", "hello")`, nil))

	// $hmac(algorithm, key, value) – hex-encoded HMAC
	fmt.Println("hmac:  ", eval(`$hmac("sha256", "secret", "message")`, nil))
}
