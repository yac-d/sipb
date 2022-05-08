function uploadFormData(formdata) {
	return new Promise(function(resolve, reject) {

		let xhr = new XMLHttpRequest();
		xhr.onreadystatechange = function() {
			if (this.readyState == 4) {
				if (this.status == 413) {
					alert("File uploaded too large! Truncated to maximum size.");
				}
				resolve("Done!");
			}
		}
		xhr.open("POST", "/upload", true);
		xhr.send(formdata);
	});
}

function refreshFileList() {
	let cntPromise = getFileCnt();
	cntPromise.then(setFileCnt);
	cntPromise.then(function() {
		let fileContainer = newFileContainer();
		let firstFileContainer = document.getElementById("container").firstElementChild;
		if (firstFileContainer == null) {
			document.getElementById("container").append(fileContainer);
		}
		else {
			firstFileContainer.insertAdjacentElement("beforebegin", fileContainer);
		}
		let filePromise = fetchLastNthUploadedFile(1);
		filePromise.then(details => populateFileContainer(details, fileContainer));
	});
}

function uploadOnClick() {
	chooser = document.getElementById("filechooser");
	if (chooser.value == "") {
		alert("File to upload not selected");
		return;
	}
	document.getElementById("spinner").style.display = "inline-block";
	let p = uploadFormData(FormData(uploadForm));
	p.then(function() {
		chooser.value = "";
		refreshFileList();
	});
	p.finally(() => document.getElementById("spinner").style.display = "none");
}

function uploadOnPaste() {
	document.getElementById("spinner").style.display = "inline-block";
	let fd = new FormData();
	fd.append("file", event.dataTransfer.files[0]);
	let p = uploadFormData(fd);
	p.then(refreshFileList);
	p.finally(() => document.getElementById("spinner").style.display = "none");
}

function handleDragAndDrop(event) {
	console.log("Dropped");
	event.preventDefault();
	console.log(event.dataTransfer);
	uploadOnPaste();
}

function handleDragOver(event) {
	event.preventDefault();
	console.log("Dragged over");
}
