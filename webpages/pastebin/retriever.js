function loadFiles() {
	let fileCntPromise = getFileCnt();
	fileCntPromise.then(setFileCnt);
	fileCntPromise.then(() => showAtMostNMoreFiles(10));
}
