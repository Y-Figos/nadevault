// static/js/uploader.js

document.addEventListener("change", function(e) {
    // Verifica se o elemento que disparou o evento tem a classe desejada
    if (e.target && e.target.classList.contains("imageUpload")) {
        const files = e.target.files;
        
        if (files.length > 0) {
            console.log("Arquivos selecionados para:", e.target.name);
            
            // Como você usou 'multiple' no seu template Templ, 
            // podemos iterar sobre todos os arquivos
            Array.from(files).forEach(file => {
                console.log("- Arquivo:", file.name, `(${file.size} bytes)`);
            });
        }
    }
});