//设置cookie
function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+d.toUTCString();
    document.cookie = cname + "=" + cvalue + "; " + expires + ";path=/";
}
//获取cookie
function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for(var i=0; i<ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1);
        if (c.indexOf(name) != -1) return c.substring(name.length, c.length);
    }
    return "";
}
//清除cookie
function clearCookie(name) {
    setCookie(name, "", -1);
}

function logout() {
    clearCookie('uid');
    location.href = "index.html";
}

// document.cookie = "uid=Faldict";
if (getCookie('uid')) {
    console.log('Welcome!');
    document.getElementById('logstate').innerHTML = "Log Out";
    document.getElementById('logstate').onclick = logout;
} else {
    console.log("please log in first!");
    document.getElementById('logstate').innerHTML = "Log In";
    document.getElementById('logstate').href = "login.html";
}
