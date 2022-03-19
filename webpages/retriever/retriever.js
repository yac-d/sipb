function showFile(details) {
	let filetype = details[0]; let path = details[1];
	let fileContainer = document.createElement("div");
	let pathElements = path.split("/");
	let name = pathElements[pathElements.length - 1];

	fileContainer.classList.add("fileContainer");

	if (filetype.includes("image")) {
		fileContainer.classList.add("imgContainer");

		let filename = document.createElement("h3");
		filename.innerText = name;

		let img = document.createElement("img");
		img.src = path;

		fileContainer.append(filename)
		fileContainer.append(img);
	}
	else {
		fileContainer.class = "fileContainer";
		let link = document.createElement("a");
		link.href = path;
		link.innerText = name;
		fileContainer.append(link);
	}
	document.getElementById("container").append(fileContainer);
}

function showLastNthUploadedFile(n) {
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
				//setFileCnt(parseInt(this.responseText));
				resolve(parseInt(this.responseText));
			}
		}

		xhr.open("GET", "/retrieve/fileCount", true);
		xhr.send();
	});
}

function showFiles(c) {
	for (i=1; i<c+1; i++) {
		p = showLastNthUploadedFile(i);
		p.then(showFile)
	}
}

let fileCntPromise = getFileCount();
fileCntPromise.then(setFileCnt);
fileCntPromise.then(showFiles);
//showLastNthUploadedFile(1);
//showLastNthUploadedFile(2);
