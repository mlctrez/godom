//
// This javascript is appended to go's wasm_exec.js
//
(() => {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
    .then((result) => {
      go.run(result.instance)
        .then(() => console.log("go.run exited"))
        .catch(err => console.log("error ", err));
    }).catch(err => console.log("error ", err))
})();
