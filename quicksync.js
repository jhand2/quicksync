var gulp = require("gulp");
var client = require("scp2");
var fs = require("fs");
var read = require("read");

// Windows IP: 137.117.54.129
// Windows OE dir: /C:/code/openenclave

var CONFIG = {
    "client_oedir": "/home/jordan/code/openenclave",
    "client_privkey": "/home/jordan/.ssh/id_rsa",
    "target_username": "jorhand",
    "target_ip": "13.68.192.102",
    "target_oedir": "/home/jorhand/openenclave"
}


function get_ssh_pass() {
    return new Promise(function(resolve) {
        read(
            { prompt: "SSH Passphrase: ", silent: true },
            function(err, pass) {
                resolve(pass);
            }
        );
    });
}

function copy_file_to_target(file, username, pass) {
    var target_file = file.replace(CONFIG.client_oedir, CONFIG.target_oedir);
    console.log(`Copying ${file} to ${target_file}`);
    client.scp(file, {
        host: CONFIG.target_ip,
        username: username,
        path: target_file,
        privateKey: fs.readFileSync(CONFIG.client_privkey),
        passphrase: pass
    }, function(err) {});
}

function set_file_watch(callback) {
    gulp.watch([`${CONFIG.client_oedir}/**/*`,
                `!${CONFIG.client_oedir}/3rdparty/**/*`,
                `!${CONFIG.client_oedir}/build/**/*`,
                `!${CONFIG.client_oedir}/.git/**/*`])
        .on("change", callback);
}

function main() {
    get_ssh_pass().then(function(pass) {
        set_file_watch(function(file) {
            copy_file_to_target(file, CONFIG.target_username, pass);
        });
    });
}

main();
