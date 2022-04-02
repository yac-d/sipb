function newFileContainer() {
	let fileContainer = document.createElement("div");
	fileContainer.classList.add("fileContainer");
	return fileContainer;
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
	let date = new Date(parseInt(name.substr(0, name.search("_"))));

	console.log(fileContainer);
	let filename = document.createElement("h3");
	let link = document.createElement("a");
	let timestamp = document.createElement("p");
	timestamp.innerText = "Uploaded " + date.toLocaleString();
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

	fileContainer.append(timestamp);
}

function setFileCnt(cnt) {
	if (cnt != 1) {
		document.getElementById("count").innerHTML = cnt.toString() + " files";
	}
	else {
		document.getElementById("count").innerHTML = "1 file";
	}
}

function getFileCnt() {
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
