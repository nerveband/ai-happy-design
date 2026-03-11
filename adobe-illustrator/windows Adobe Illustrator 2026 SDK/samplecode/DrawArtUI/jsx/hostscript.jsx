/*
 * Extendscript scripts should be located here.
 */
function checkDocumentIsOpen() {
    if (app.documents.length === 0) {
        alert("A document is required to load a swatch");
    }
}
