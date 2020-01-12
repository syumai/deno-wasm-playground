import './wasm_exec.js';

const go = new Go();

const result = await WebAssembly.instantiate(
  Deno.readFileSync('./main.wasm'),
  go.importObject
);

await go.run(result.instance);
