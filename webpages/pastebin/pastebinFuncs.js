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
				resolve(resp);
			}
		};

		xhr.open("POST", "/retrieve", true);
		xhr.send(n.toString());
	});
}

function logBase(b, e) {
	return Math.log(e) / Math.log(b);
}

function prettySize(bytes) {
	let suffixes = ["B", "KiB", "MiB", "GiB"];
	let logB1024 = Number.parseInt(logBase(1024, bytes));
	let suffix = suffixes[logB1024];
	let num = (bytes / (1024 ** logB1024)).toFixed(2);
	return num.toString() + " " + suffix;
}

function populateFileContainer(details, fileContainer) {
    let filelocation = "/static/" + details.ID;

	let filename = document.createElement("h3");
	let link = document.createElement("a");
	filename.innerText = details.Name;
	link.classList.add("filelink");
    link.href = filelocation;
    link.download = details.Name
	link.append(filename);
	fileContainer.append(link);

	let fileInfoBox = document.createElement("div");
	fileInfoBox.classList.add("fileinfobox");

	let sizestamp = document.createElement("p");
	sizestamp.innerText = prettySize(details.Size);
	sizestamp.classList.add("sizestamp");
	fileInfoBox.append(sizestamp);

	let timestamp = document.createElement("p");
    let options = {
        weekday: 'long', year: 'numeric', month: 'long', day: 'numeric', 
        timeZone: "Asia/Kolkata", hour: 'numeric', minute: 'numeric', hour12: false
    };
	timestamp.innerText = Intl.DateTimeFormat("default", options).format(Date.parse(details.Timestamp));
	timestamp.classList.add("timestamp");
	fileInfoBox.append(timestamp);

	fileContainer.append(fileInfoBox);

	if (details.Type.includes("image")) {
		fileContainer.classList.add("imgContainer");

		let img = document.createElement("img");
		img.src = filelocation;
		img.classList.add("imagePreview");

		fileContainer.append(img);
	}

    if (details.Note != "") {
        let note = document.createElement("p");
        note.classList.add("note");
        note.innerText = details.Note;
        fileContainer.append(note)
    }
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

function showAtMostNMoreFiles(n) {
	let alreadyShownCnt = document.getElementById("container").children.length;
	let available = parseInt(document.getElementById("count").innerText);
	if (alreadyShownCnt == available) return;
	for (i=alreadyShownCnt+1; i<=Math.min(alreadyShownCnt+n, available); i++) {
		let fileContainer = newFileContainer();
		document.getElementById("container").append(fileContainer);

		p = fetchLastNthUploadedFile(i);
		p.then(details => populateFileContainer(details, fileContainer));
	}
}
