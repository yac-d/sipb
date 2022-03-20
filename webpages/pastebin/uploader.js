function uploadFile(form) {
	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/upload", true);
	xhr.send(new FormData(uploadForm));
}
