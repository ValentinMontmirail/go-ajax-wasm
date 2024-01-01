// Instantiate the Go object.
const go = new Go();

// Override the 'syscall/js.finalizeRef' method with a mock function, since 'finalizeRef' is not functional yet.
go.importObject.gojs["syscall/js.finalizeRef"] = _ => 0;  // 😉 Hope to remove it one day...

function loadScript(url) {
    return new Promise((resolve, reject) => {
        const script = document.createElement('script');
        script.src = url;
        script.onload = resolve;
        script.onerror = reject;
        document.body.appendChild(script);
    });
}

WebAssembly.instantiateStreaming(fetch("authors.wasm"), go.importObject).then((result) => {
    
    go.run(result.instance);
    
    // After Wasm is loaded, load the other scripts
    loadScript("js/authors.js").then(() => {
        console.log("Authors.js loaded");
        return loadScript("js/index.js");
    }).then(() => {
        console.log("Index.js loaded");
    }).catch(err => {
        console.error("Error loading scripts:", err);
    });
});
