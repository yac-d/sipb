function uploadFile() {
	return new Promise(function(resolve, reject) {
		chooser = document.getElementById("filechooser");
		if (chooser.value == "") {
			reject(new Error("File to upload not selected"));
			return;
		}

		let xhr = new XMLHttpRequest();
		xhr.onreadystatechange = function() {
			if (this.readyState == 4) {
				resolve("Done!");
			}
		}
		xhr.open("POST", "/upload", true);
		xhr.send(new FormData(uploadForm));
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

function uploadAndRefresh() {
	document.getElementById("spinner").style.display = "inline-block";
	let p = uploadFile();
	p.then(function() {
		chooser.value = "";
		refreshFileList();
	});
	p.catch((e) => alert("Select a file first!"));
	p.finally(() => document.getElementById("spinner").style.display = "none");
}
