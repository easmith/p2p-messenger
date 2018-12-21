import React, {Component} from 'react';
import './App.css';

class App extends Component {


    constructor(props, context) {
        super(props, context);

        this.state = {
            socket: null,
            messages: ""
        }


    }

    componentDidMount ()  {

        let socket = new WebSocket("ws://localhost:35035/ws");

        this.setState({socket: socket});

        socket.onopen = function() {
            console.log("Соединение установлено.");
            socket.send("hello from front")
        };

        socket.onmessage = this.onMessage;

        socket.onclose = function(event) {
            if (event.wasClean) {
                console.log('Соединение закрыто чисто');
            } else {
                console.log('Обрыв соединения');
            }
            console.log('Код: ' + event.code + ' причина: ' + event.reason);
        };

        socket.onerror = function(error) {
            console.log("Ошибка " + error.message);
        };

    }

    onMessage = (event) => {
        this.setState({
            messages:  this.state.messages + "\n" + event.data
        });
        console.log("Получены данные " + event.data);
    }

    render() {
        return (
            <div className="App">
                P2P-messenger - training
                <pre>{this.state.messages}</pre>
                <input type={"text"} onKeyPress={this._handleEnter}/>
            </div>
        );
    }

    _handleEnter = (e) => {
        if (e.key === 'Enter') {
            this.state.socket.send(e.target.value);
            e.target.value = "";
        }
    }


}

export default App;
