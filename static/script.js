document.getElementById("form-id").onsubmit = function () {
    var urlInput = document.getElementById("inputURL").value;
    if (isValidURL(urlInput)) {
        var id = getIdFromURL(urlInput);
        document.getElementById("inputID").value = id;
    } else {
        alert("Invalid URL. Please use a URL in the format: https://sankaku.app/post/show/12345678");
        return false;
    }
};

function isValidURL(url) {
    var regex = /^https:\/\/[^/]+\/post\/show\/\d+(\?.*)?$/;
    return regex.test(url);
}

function getIdFromURL(url) {
    var cleanURL = url.split('?')[0];
    var parts = cleanURL.split('/');
    var idPart = parts[parts.length - 1];
    return idPart;
}
