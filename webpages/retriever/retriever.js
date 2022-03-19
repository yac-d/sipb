function showFile(filetype, path) {
	if (filetype.includes("image")) {
		let img = document.createElement("img");
		img.src = path;
		document.getElementById("container").append(img);
	}
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

	xhr.open("POST", "/retrieve", false);
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

	xhr.open("GET", "/retrieve/fileCount", false);
	xhr.send();
}

getFileCount();
showLastNthUploadedFile(1);
