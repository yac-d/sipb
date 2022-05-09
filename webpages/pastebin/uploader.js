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
	let p = uploadFormData(new FormData(uploadForm));
	p.then(function() {
		chooser.value = "";
		refreshFileList();
	});
	p.finally(() => document.getElementById("spinner").style.display = "none");
}

function uploadOnDrop(event) {
	event.preventDefault();
	let file = event.dataTransfer.files[0];
	if (file.size < 1) { // Folders have 0 size
		alert("Folder or empty file. Ignoring.");
		handleDragLeave();
		return;
	}

	document.getElementById("spinner").style.display = "inline-block";
	let fd = new FormData();
	fd.append("file", file);

	let p = uploadFormData(fd);
	p.then(refreshFileList);
	p.finally(() => document.getElementById("spinner").style.display = "none");

	handleDragLeave();
}

function handleDragover() {
	event.preventDefault();
	document.getElementById("dropzone").classList.add("dropzoneHover");
}
function handleDragLeave() {
	document.getElementById("dropzone").classList.remove("dropzoneHover");
}

document.addEventListener("drop", uploadOnDrop);
