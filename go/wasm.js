import './wasm_exec.js';

(async () => {
  const go = new Go();
  const result = await WebAssembly.instantiate(
    Deno.readFileSync('./main.wasm'),
    go.importObject
  );
  const mod = result.module;
  let inst = result.instance;

  console.clear();
  await go.run(inst);
  inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
})();
