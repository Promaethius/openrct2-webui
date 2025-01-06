/// <reference path="./openrct2.d.ts" />

function main() {
    if (network.mode === 'server') {

        let server = network.createListener();
        let port = context.sharedStorage.get('remote.port', 35711);
        let host = context.sharedStorage.get('remote.host', '127.0.0.1');

        server.on('connection', (socket) => {
            socket.on('data', (data) => {
                socket.write(JSON.stringify(eval(data)));
            });
        });

        server.listen(port, host);
    }
}

registerPlugin({
    name: 'remote',
    version: '0.0.1',
    authors: ['Jonathan Bryant'],
    type: 'remote',
    minApiVersion: 19,
    targetApiVersion: 83,
    licence: 'MIT',
    main
});