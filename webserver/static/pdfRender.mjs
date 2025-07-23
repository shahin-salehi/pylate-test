// If absolute URL from the remote server is provided, configure the CORS

// header on that server.





// Loaded via <script> tag, create shortcut to access PDF.js exports.

//var { pdfjsLib } = globalThis;



// The workerSrc property shall be specified.

pdfjsLib.GlobalWorkerOptions.workerSrc = '/public/node_modules/pdfjs-dist/build/pdf.worker.mjs';





// data from server

// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/script

const instructions = JSON.parse(document.getElementById("viewer-instructions").getAttribute("data"));

console.log("instructions: %o", instructions)



/**

 * LETS GET THIS VIEWER UP AND RUNNING

 *   https://github.com/mozilla/pdf.js/blob/cf5a1d60a61e13108afffbb1046521b5bb46f918/examples/components/simpleviewer.js#L38

 */



const CMAP_URL = "/public/node_modules/pdfjs-dist/cmaps/";

const CMAP_PACKED = true;







const DEFAULT_URL = instructions.path;



const ENABLE_XFA = true;

const SEARCH_FOR = instructions.highlight;



const SANDBOX_BUNDLE_SRC = new URL(

    "/public/node_modules/pdfjs-dist/build/pdf.sandbox.mjs",

    window.location

);

const container = document.getElementById("viewerContainer");



const eventBus = new pdfjsViewer.EventBus();



// (Optionally) enable hyperlinks within PDF files.

const pdfLinkService = new pdfjsViewer.PDFLinkService({

    eventBus,

});



// (Optionally) enable find controller.

const pdfFindController = new pdfjsViewer.PDFFindController({

    eventBus,

    linkService: pdfLinkService,

});



// (Optionally) enable scripting support.

const pdfScriptingManager = new pdfjsViewer.PDFScriptingManager({

    eventBus,

    sandboxBundleSrc: SANDBOX_BUNDLE_SRC,

});





const pdfViewer = new pdfjsViewer.PDFViewer({

    container,

    eventBus,

    linkService: pdfLinkService,

    findController: pdfFindController,

    scriptingManager: pdfScriptingManager,

});

pdfLinkService.setViewer(pdfViewer);

pdfScriptingManager.setViewer(pdfViewer);





eventBus.on("pagesinit", function () {

    // We can use pdfViewer now, e.g. let's change default scale.

    console.log("pdfviewerobject: ", pdfViewer)

    //pdfViewer.currentScaleValue = "page-fit";



    // We can try searching for things.

    if (SEARCH_FOR) {

        eventBus.dispatch("find", { type: "", query: SEARCH_FOR });

    }

});





// Asynchronous download of PDF

// Loading document.

const loadingTask = pdfjsLib.getDocument({

    url: DEFAULT_URL,

    cMapUrl: CMAP_URL,

    cMapPacked: CMAP_PACKED,

    enableXfa: ENABLE_XFA,

});





const pdfDocument = await loadingTask.promise;

// Document loaded, specifying document for the viewer and

// the (optional) linkService.

pdfViewer.setDocument(pdfDocument);



pdfLinkService.setDocument(pdfDocument, null);

