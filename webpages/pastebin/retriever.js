function showFiles(c) {
	for (i=1; i<c+1; i++) {
		let fileContainer = newFileContainer();
		document.getElementById("container").append(fileContainer);

		p = fetchLastNthUploadedFile(i);
		p.then(details => populateFileContainer(details, fileContainer));
	}
}

function loadFiles() {
	document.getElementById("container").innerHTML = "";
	let fileCntPromise = getFileCnt();
	fileCntPromise.then(setFileCnt);
	fileCntPromise.then(showFiles);
}
