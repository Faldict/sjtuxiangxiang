// ---------------------------------------
//
// This is the user interface javascript.
// Author : Faldict
// Date: 2017-06-01
//
// ---------------------------------------

//设置cookie
function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+d.toUTCString();
    document.cookie = cname + "=" + cvalue + "; " + expires +";path=/";
}

function login() {
    var username = document.getElementById('username').value;
    var passwd = document.getElementById('passwd').value;
    var uri = "/user/login";
    $.post(uri, {
        "username" : username,
        "password" : passwd
    }, function (status) {
        if (status == "200000") {
            setCookie('uid', username, 3);
            console.log("Welcome!");
            location.href = "index.html";
        } else {
            alert("Login Error!");
            console.log("Error Code: " + status + '\n');
        }
    });
}

function register() {
    var username = document.getElementById('username').value;
    var passwd = document.getElementById('passwd').value;
    var email = document.getElementById('email').value;
    var uri = "/user/register";
    $.post(uri, {
        "username" : username,
        "password" : passwd,
        "email" : email
    }, function (info) {
        alert(info);
    });
    window.location = "login.html";
}
