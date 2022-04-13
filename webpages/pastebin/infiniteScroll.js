let hasScrolled = false;

function loadMore() {
	if (window.scrollY > document.body.scrollHeight*0.8 && !hasScrolled) {
		hasScrolled = true;
		showAtMostNMoreFiles(5)
		setTimeout(() => {hasScrolled = false}, 1000);
	}
}

window.onscroll = loadMore;
