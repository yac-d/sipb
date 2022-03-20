function showFile(details) {
	let filetype = details[0]; let path = details[1];
	let fileContainer = document.createElement("div");
	let pathElements = path.split("/");
	let name = pathElements[pathElements.length - 1];
	let prettyName = name.substr(name.search("_")+1, name.length);

	fileContainer.classList.add("fileContainer");
	let filename = document.createElement("h3");
	let link = document.createElement("a");
	link.href = path;
	filename.innerText = prettyName;
	link.append(filename);
	fileContainer.append(link);

	if (filetype.includes("image")) {
		fileContainer.classList.add("imgContainer");

		let img = document.createElement("img");
		img.src = path;
		img.classList.add("imagePreview");

		fileContainer.append(img);
	}
	document.getElementById("container").append(fileContainer);
}

function fetchLastNthUploadedFile(n) {
	return new Promise(function(resolve, reject) {
		let xhr = new XMLHttpRequest();

		xhr.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
				let resp = {};
				resp = JSON.parse(this.responseText);
				console.log(resp);
				resolve([resp.Type, resp.Path]);
			}
		};

		xhr.open("POST", "/retrieve", true);
		xhr.send(n.toString());
	});
}

function setFileCnt(cnt) {
	document.getElementById("count").innerHTML = cnt.toString() + " files";
	console.log(cnt);
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

function showFiles(c) {
	for (i=1; i<c+1; i++) {
		p = fetchLastNthUploadedFile(i);
		p.then(showFile)
	}
}

function loadFiles() {
	let fileCntPromise = getFileCount();
	fileCntPromise.then(setFileCnt);
	fileCntPromise.then(showFiles);
}
