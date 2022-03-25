function uploadFile(form) {
	chooser = document.getElementById("filechooser");
	if (chooser.value == "") {
		alert("Select a file first!");
		return;
	}
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/upload", true);
	xhr.send(new FormData(uploadForm));
	chooser.value = "";
	loadFiles();
}
