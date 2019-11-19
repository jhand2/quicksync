var gulp = require("gulp");
var client = require("scp2");
var fs = require("fs");
var read = require("read");
var git = require("nodegit")

var args = process.argv.slice(2)
if (args.length < 1) {
    console.log("Usage: node quicksync.js <config_name>");
    process.exit(1);
}

var CONFIG = read_config(args[0]);
console.log(CONFIG);

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

function read_config(config_name) {
    var f = fs.readFileSync("/home/jordan/.quicksync.conf");
    return JSON.parse(f)[config_name];
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
        // Copy any files that are not currently checked into git
        git.Repository.open(CONFIG.client_oedir).then(function(repo) {
            return repo.getStatus();
        }).then(function(results) {
            for (let r of results) {
                var f = CONFIG.client_oedir + "/" + r.path();
                copy_file_to_target(f, CONFIG.target_username, pass);
            }
        });

        // Watch for future file changes
        set_file_watch(function(file) {
            copy_file_to_target(file, CONFIG.target_username, pass);
        });
    });
}

main();
