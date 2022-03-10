var xhr = new XMLHttpRequest();

xhr.onreadystatechange = function() {
	if (this.readyState == 4 && this.status == 200) {
		console.log("Request complete.");
	}
};

xhr.open("POST", "/retrieve", true);
xhr.send("5");
