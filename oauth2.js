var system = require('system');
var env = system.env;
var page = require('webpage').create();

page.onConsoleMessage = function() {};
page.onError = function() {};

var loginUrl = env["CONCOURSE_URL"] + "/auth/github?team_name=" + env["CONCOURSE_TARGET"];

var userName = env["GITHUB_LDAP_USERNAME"];
var password = env["GITHUB_LDAP_PASSWORD"];

var bearer = "Bearer ";
//console.log(userName,password);

page.open(loginUrl, function() {

    page.evaluate(function(user,password) {
        document.querySelector("input[name='login']").value = user;
        document.querySelector("input[name='password']").value = password;
        document.forms[0].submit();
    },userName,password);

    window.setTimeout(function() {
        var cookies = page.cookies;

        for(var i in cookies) {
            if (cookies[i].value.substr(0, bearer.length) == bearer) {
                console.log(cookies[0].value);
                phantom.exit();
                return
            }
        }
        phantom.exit(1);
    },5000);

});
