// ---------------------------------------
//
// This is the user interface javascript.
// Author : Faldict
// Date: 2017-06-01
//
// ---------------------------------------

function login() {
    var username = document.getElementById('username').value;
    var passwd = document.getElementById('passwd').value;
    var uri = "/user/login";
    $.post(uri, {
        "username" : username,
        "password" : passwd
    }, function (status) {
        if (status == "200000") {
            alert("Login Success!");
        } else {
            alert("Login Error!");
            console.log("Error Code: " + status + '\n');
        }
    });
    window.location = "index.html";
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