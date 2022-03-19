function showFile(filetype, path) {
	let fileContainer = document.createElement("div");
	let pathElements = path.split("/");
	let name = pathElements[pathElements.length - 1];

	if (filetype.includes("image")) {
		fileContainer.class = "fileContainer imgContainer";

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
	let xhr = new XMLHttpRequest();

	xhr.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			let resp = {};
			resp = JSON.parse(this.responseText);
			console.log(resp);
			showFile(resp.Type, resp.Path);
		}
	};

	xhr.open("POST", "/retrieve", true);
	xhr.send(n.toString());
}

function showFileCnt(cnt) {
	document.getElementById("count").innerHTML = cnt.toString() + " files";
	console.log(cnt);
}

function getFileCount() {
	let xhr = new XMLHttpRequest();

	xhr.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			showFileCnt(parseInt(this.responseText));
		}
	}

	xhr.open("GET", "/retrieve/fileCount", true);
	xhr.send();
}

getFileCount();
showLastNthUploadedFile(1);
showLastNthUploadedFile(2);
