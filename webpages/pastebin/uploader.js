function uploadFile() {
	return new Promise(function(resolve, reject) {
		chooser = document.getElementById("filechooser");
		if (chooser.value == "") {
			reject(new Error("File to upload not selected"));
			return;
		}

		let xhr = new XMLHttpRequest
		xhr.onreadystatechange = function() {
			if (this.readyState == 4) {
				console.log("Request complete.");
				resolve("Done!");
			}
		}
		xhr.open("POST", "/upload", true);
		xhr.send(new FormData(uploadForm));
	});
}

function uploadAndRefresh() {
	let p = uploadFile();
	p.then(function() {
		chooser.value = "";
		loadFiles();
	});
	p.catch((e) => alert("Select a file first!"));
}
