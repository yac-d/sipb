function setFileCnt(cnt) {
	if (cnt != 1) {
		document.getElementById("count").innerHTML = cnt.toString() + " files";
	}
	else {
		document.getElementById("count").innerHTML = "1 file";
	}
}

function getFileCount() {
	return new Promise(function(resolve, reject) {
		let xhr = new XMLHttpRequest();

		xhr.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
				resolve(parseInt(this.responseText));
			}
		}

		xhr.open("GET", "/retrieve/fileCount", true);
		xhr.send();
	});
}

function fetchLastNthUploadedFile(n) {
	return new Promise(function(resolve, reject) {
		let xhr = new XMLHttpRequest();

		xhr.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
				let resp = {};
				resp = JSON.parse(this.responseText);
				console.log(resp);
				resolve(resp);
			}
		};

		xhr.open("POST", "/retrieve", true);
		xhr.send(n.toString());
	});
}

function populateFileContainer(details, fileContainer) {
	let pathElements = details.Path.split("/");
	let name = pathElements[pathElements.length - 1];
	let prettyName = name.substr(name.search("_")+1, name.length);

	console.log(fileContainer);
	let filename = document.createElement("h3");
	let link = document.createElement("a");
	link.href = details.Path;
	filename.innerText = prettyName;
	link.append(filename);
	fileContainer.append(link);

	if (details.Type.includes("image")) {
		fileContainer.classList.add("imgContainer");

		let img = document.createElement("img");
		img.src = details.Path;
		img.classList.add("imagePreview");

		fileContainer.append(img);
	}
}

function showFiles(c) {
	for (i=1; i<c+1; i++) {
		let fileContainer = document.createElement("div");
		fileContainer.id = "fileContainer" + i.toString();
		fileContainer.classList.add("fileContainer");
		document.getElementById("container").append(fileContainer);

		p = fetchLastNthUploadedFile(i);
		p.then(details => populateFileContainer(details, fileContainer));
	}
}

function loadFiles() {
	document.getElementById("container").innerHTML = "";
	let fileCntPromise = getFileCount();
	fileCntPromise.then(setFileCnt);
	fileCntPromise.then(showFiles);
}
