import React, {Component} from 'react';
import './App.css';

class App extends Component {

    componentDidMount() {

        var socket = new WebSocket("ws://localhost:35035/ws");

        socket.onopen = function() {
            console.log("Соединение установлено.");
        };

        socket.onclose = function(event) {
            if (event.wasClean) {
                console.log('Соединение закрыто чисто');
            } else {
                console.log('Обрыв соединения'); // например, "убит" процесс сервера
            }
            console.log('Код: ' + event.code + ' причина: ' + event.reason);
        };

        socket.onmessage = function(event) {
            console.log("Получены данные " + event.data);
        };

        socket.onerror = function(error) {
            console.log("Ошибка " + error.message);
        };

    }

    render() {
        return (
            <div className="App">
                P2P-messenger - training
            </div>
        );
    }
}

export default App;
