import './wasm_exec.js';

if (!WebAssembly.instantiateStreaming) {
  // polyfill
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

(async () => {
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch('./test.wasm'),
    go.importObject,
  );
  const mod = result.module;
  const inst = result.instance;

  console.clear();
  await go.run(inst);
  inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
})();
