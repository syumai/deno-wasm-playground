import mainwasm from "./mainwasm.ts";
import { Go } from "./wasm_exec.js";
import { decode } from "https://deno.land/std@0.92.0/encoding/base64.ts";
import { Buffer } from "https://deno.land/std@0.92.0/io/buffer.ts";

const bytes = decode(mainwasm);
const go = new Go();
const result = await WebAssembly.instantiate(bytes, go.importObject);
go.run(result.instance);

async function readExample() {
    const f = await Deno.open("./example.txt");
    const result = read(f);
    console.log(result);
    f.close();
}

async function writeExample() {
    const buf = new Buffer();
    const result = write(buf);
    console.log(result);
}

async function seekExample() {
    const f = await Deno.open("./example.txt");
    const result = seek(f);
    console.log(result);
    f.close();
}

await seekExample();
await writeExample();
await readExample();
callExit();
