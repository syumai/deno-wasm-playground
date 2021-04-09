//go:generate sh -c "GOOS=js GOARCH=wasm go build -o main.wasm ./ && cat main.wasm | deno run https://denopkg.com/syumai/binpack/mod.ts > mainwasm.ts && rm main.wasm"
package main

import (
	"io"
	"syscall/js"
)

func main() {
	var r js.Func
	r = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer r.Release()
		var cb js.Func
		cb = js.FuncOf(func (_ js.Value, pArgs []js.Value) interface{} {
			defer cb.Release()
			resolve := pArgs[0]
			b := make([]byte, 16)
			go func() {
				_, err := read(args[0], b) // read 16 bytes from JS
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.ValueOf(string(b)))
			}()
			return js.Undefined()
		})
		p := newPromise(cb)
		return p
	})
	js.Global().Set("read", r)

	var w js.Func
	w = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer w.Release()
		var cb js.Func
		cb = js.FuncOf(func (_ js.Value, pArgs []js.Value) interface{} {
			defer cb.Release()
			resolve := pArgs[0]
			b := []byte{1, 2, 3, 4, 5, 6, 7, 8}
			go func() {
				n, err := write(args[0], b) // write 8 bytes to JS
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.ValueOf(n))
			}()
			return js.Undefined()
		})
		p := newPromise(cb)
		return p
	})
	js.Global().Set("write", w)

	var s js.Func
	s = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer s.Release()
		var cb js.Func
		cb = js.FuncOf(func (_ js.Value, pArgs []js.Value) interface{} {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				n, err := seek(args[0], 4, 0) // seek 4 bytes from start
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.ValueOf(n))
			}()
			return js.Undefined()
		})
		p := newPromise(cb)
		return p
	})
	js.Global().Set("seek", s)
	select {}
}

func read(v js.Value, p []byte) (int, error) {
	ua := newUint8Array(len(p))
	promise := v.Call("read", ua)
	resultCh := make(chan js.Value)
	eofCh := make(chan struct{})

	var then, catch js.Func
	then = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer then.Release()
		result := args[0]
		if result.IsNull() {
			eofCh <- struct{}{}
		}
		resultCh <- result
		return js.Undefined()
	})
	catch = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer catch.Release()
		result := args[0]
		close(resultCh)
		panic(result)
	})
	promise.Call("then", then).Call("catch", catch)
	select {
	case result := <-resultCh:
		_ = js.CopyBytesToGo(p, ua)
		return result.Int(), nil
	case <-eofCh:
		return 0, io.EOF
	}
}

func write(v js.Value, p []byte) (int, error) {
	ua := newUint8Array(len(p))
	_ = js.CopyBytesToJS(ua, p)
	promise := v.Call("write", ua)
	resultCh := make(chan js.Value)

	var then, catch js.Func
	then = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer then.Release()
		resultCh <- args[0]
		return js.Undefined()
	})
	catch = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer catch.Release()
		close(resultCh)
		panic(args[0])
	})
	promise.Call("then", then).Call("catch", catch)
	result := <-resultCh
	return result.Int(), nil
}

func seek(v js.Value, offset int64, whence int) (int64, error) {
	resultCh := make(chan js.Value)
	go func() {
		promise := v.Call("seek", js.ValueOf(offset), js.ValueOf(whence))

		var then, catch js.Func
		then = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			defer then.Release()
			resultCh <- args[0]
			return js.Undefined()
		})
		catch = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			defer catch.Release()
			close(resultCh)
			panic(args[0])
		})
		promise.Call("then", then).Call("catch", catch)

	}()
	result := <-resultCh
	return int64(result.Int()), nil
}

type Callback = func(this js.Value, args []js.Value) interface{}

func newPromise(fn js.Func) js.Value {
	p := js.Global().Get("Promise")
	return p.New(fn)
}

func newUint8Array(size int) js.Value {
	ua := js.Global().Get("Uint8Array")
	return ua.New(size)
}
