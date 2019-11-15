var gulp = require("gulp");
var client = require("scp2");
var fs = require("fs");
var read = require("read");

// Dev machine files
var HOME = "/home/jordan";
var oe_dir = `${HOME}/code/openenclave`
var ssh_key = `${HOME}/.ssh/id_rsa`;

// Target
var user = "jorhand"
var ip = "137.117.54.129";
var oe_target = '/C:/code/openenclave';

read(
    { prompt: "SSH Passphrase: ", silent: true },
    function(err, pass) {
        set_file_watch(function(file) {
            copy_file_to_target(file, user, pass);
        })
    }
)

function copy_file_to_target(file, username, pass) {
    var target_file = file.replace(oe_dir, oe_target);
    console.log(`Copying ${file} to ${target_file}`);
    client.scp(file, {
        host: ip,
        username: username,
        path: target_file,
        privateKey: fs.readFileSync(ssh_key),
        passphrase: pass
    }, function(err) {});
}

function set_file_watch(callback) {
    gulp.watch([`${oe_dir}/**/*`,
                `!${oe_dir}/3rdparty/**/*`,
                `!${oe_dir}/build/**/*`,
                `!${oe_dir}/.git/**/*`])
        .on("change", callback);
}

