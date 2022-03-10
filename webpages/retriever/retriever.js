function showFile(filetype, path) {
	if (filetype.includes("image")) {
		document.getElementById("img1").src = path;
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

	xhr.open("POST", "/retrieve", true);
	xhr.send(n.toString());
}

showLastNthUploadedFile(1);
