import express from 'express';
import { createServer } from 'http';
import { Server } from 'socket.io';
import { startAttack, stopAttack } from './attack';

const app = express();
const server = createServer(app);
const io = new Server(server);

app.use(express.static('public'));

io.on('connection', (socket) => {
    console.log('LDS: a user connected');

    socket.on('startAttack', (targetUrl: string) => {
        startAttack(targetUrl, (message: string) => {
            io.emit('log', message);
        });
    });

    socket.on('stopAttack', () => {
        stopAttack();
        io.emit('log', 'LDS: Attack stopped.');
    });

    socket.on('disconnect', () => {
        console.log('LDS: user disconnected');
    });
});

server.listen(5252, () => {
    console.log('LDS: Server is running on http://localhost:5252');
});
