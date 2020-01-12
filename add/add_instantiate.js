(async () => {
  const wasm = await WebAssembly.instantiate(Deno.readFileSync('./add.wasm'));

  const { add } = wasm.instance.exports;

  console.log(add(1, 2));
})();
